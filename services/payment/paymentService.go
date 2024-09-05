package paymentService

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"time"
)

func initializePayment(req models.InitializePaymentRequest, db *storage.Database) (models.InitializePaymentResponse, int, error) {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")
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
		return models.InitializePaymentResponse{}, http.StatusBadRequest, err
	}

	payment := models.Payments{
		ID:            paystackResponse.Data.Reference,
		TicketID:      req.TicketID,
		PaymentMethod: req.PaymentMethod,
		Amount:        req.Amount,
		PaymentStatus: 1,
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
		PaymentStatus: 1,
		PaymentTime:   time.Now().Format(time.RFC3339),
	}, http.StatusOK, nil

}
