package ticketService

import (
	"errors"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/utils"
	"time"
)

func CreateTicket(req models.CreateTicketRequest, db *storage.Database) (models.CreateTicketResponse, int, error) {
	//check if ticket exists
	if ok := postgresql.CheckExistsInTable(db.Postgresql, "tickets", "event_id = ? AND user_id = ?", req.EventID, req.UserID); ok {
		return models.CreateTicketResponse{}, http.StatusConflict, errors.New("ticket for this event and user already exists")
	}
	// create ticket
	ticket := models.Ticket{}
	var responseData models.CreateTicketResponse
	ticket = models.Ticket{
		ID:           utils.GenerateUUID(),
		EventID:      req.EventID,
		UserID:       req.UserID,
		QRCode:       utils.GenerateUUID(),
		PurchaseTime: time.Now(),
		Status:       models.TicketStatusPending,
		Amount:       req.Amount,
		Type:         req.Type,
	}
	var event models.Event
	event.ID = req.EventID
	err := event.GetEventByID(db.Postgresql)
	if err != nil {
		return responseData, http.StatusNotFound, err
	}

	if int(event.TicketCount)+int(ticket.Amount) > event.Capacity {
		return responseData, http.StatusConflict, errors.New("event is full")
	}
	// Start a transaction
	tx := db.Postgresql.Begin()

	// Create ticket
	err = ticket.CreateTicket(tx)
	if err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return responseData, http.StatusInternalServerError, err
	}

	// Increase ticket count
	err = event.IncreaseTicketCount(tx, int(ticket.Amount))
	if err != nil {
		// Rollback the transaction if there's an error
		tx.Rollback()
		return responseData, http.StatusInternalServerError, err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return models.CreateTicketResponse{}, http.StatusInternalServerError, err
	}

	responseData = models.CreateTicketResponse{
		ID:           ticket.ID,
		EventID:      ticket.EventID,
		UserID:       ticket.UserID,
		QRCode:       ticket.QRCode,
		PurchaseTime: ticket.PurchaseTime,
		Status:       ticket.Status,
		Amount:       ticket.Amount,
		Type:         ticket.Type,
	}

	return responseData, http.StatusCreated, nil

}
