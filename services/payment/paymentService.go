package paymentService

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"time"
)

func InitializePayment(req models.InitializePaymentRequest, db *storage.Database) (models.InitializePaymentResponse, int, error) {
	var ticket models.Ticket
	ticket.ID = req.TicketID
	err := ticket.GetTicketByID(db.Postgresql)
	if err != nil {
		return models.InitializePaymentResponse{}, http.StatusNotFound, err
	}

	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")
	if secretKey == "" {
		return models.InitializePaymentResponse{}, http.StatusBadRequest, errors.New("PAYSTACK_SECRET_KEY is not configured")
	}

	paystackReq := models.PaystackInitializeRequest{
		Reference: "ts" + req.TicketID,
		Amount:    int(req.Amount * 100),
		Email:     req.Email,
		Currency:  "NGN",
	}
	jsonRequest, err := json.Marshal(paystackReq)
	if err != nil {
		return models.InitializePaymentResponse{}, http.StatusInternalServerError, err
	}
	// Initialize payment
	url := "https://api.paystack.co/transaction/initialize"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return models.InitializePaymentResponse{}, http.StatusInternalServerError, err
	}
	request.Header.Set("Authorization", "Bearer "+secretKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return models.InitializePaymentResponse{}, http.StatusInternalServerError, err
	}
	defer response.Body.Close()

	var paystackResponse models.PaystackInitializeResponse
	err = json.NewDecoder(response.Body).Decode(&paystackResponse)
	if err != nil {
		return models.InitializePaymentResponse{}, http.StatusInternalServerError, err
	}

	if !paystackResponse.Status {
		return models.InitializePaymentResponse{}, http.StatusBadRequest, errors.New(paystackResponse.Message)
	}

	payment := models.Payments{
		ID:            paystackResponse.Data.Reference,
		TicketID:      req.TicketID,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		PaymentStatus: models.PendingPayment,
		PaymentTime:   time.Now(),
	}

	err = payment.Create(db.Postgresql)
	if err != nil {
		return models.InitializePaymentResponse{}, http.StatusInternalServerError, err
	}

	return models.InitializePaymentResponse{
		ID:            paystackResponse.Data.Reference,
		TicketID:      req.TicketID,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		PaymentStatus: models.PendingPayment,
		PaymentTime:   time.Now().Format(time.RFC3339),
	}, http.StatusOK, nil

}

func GetPaymentByID(paymentID string, db *storage.Database) (*models.Payments, int, error) {
	var payment models.Payments
	res, err := payment.GetPaymentByID(db.Postgresql, paymentID)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return res, http.StatusOK, nil
}

func UpdatePaymentStatus(paymentID string, req models.UpdatePaymentStatusRequest, db *storage.Database) (models.UpdatePaymentStatusResponse, int, error) {
	if req.Status < models.SuccessfulPayment || req.Status > models.FailedPayment {
		return models.UpdatePaymentStatusResponse{}, http.StatusBadRequest, errors.New("invalid payment status")
	}

	var payment models.Payments
	currentPayment, err := payment.GetPaymentByID(db.Postgresql, paymentID)
	if err != nil {
		return models.UpdatePaymentStatusResponse{}, http.StatusNotFound, err
	}

	var ticket models.Ticket
	ticket.ID = currentPayment.TicketID
	err = ticket.GetTicketByID(db.Postgresql)
	if err != nil {
		return models.UpdatePaymentStatusResponse{}, http.StatusNotFound, err
	}

	ticketStatus := models.TicketStatusPending
	if req.Status == models.SuccessfulPayment {
		ticketStatus = models.TicketStatusApproved
	}
	if req.Status == models.FailedPayment {
		ticketStatus = models.TicketStatusCancelled
	}

	tx := db.Postgresql.Begin()

	err = payment.UpdatePaymentStatus(tx, paymentID, req.Status)
	if err != nil {
		tx.Rollback()
		return models.UpdatePaymentStatusResponse{}, http.StatusInternalServerError, err
	}

	err = ticket.UpdateStatus(tx, ticketStatus)
	if err != nil {
		tx.Rollback()
		return models.UpdatePaymentStatusResponse{}, http.StatusInternalServerError, err
	}

	if err := tx.Commit().Error; err != nil {
		return models.UpdatePaymentStatusResponse{}, http.StatusInternalServerError, err
	}

	return models.UpdatePaymentStatusResponse{
		PaymentID:     currentPayment.ID,
		TicketID:      currentPayment.TicketID,
		PaymentStatus: req.Status,
		TicketStatus:  ticketStatus,
	}, http.StatusOK, nil
}
