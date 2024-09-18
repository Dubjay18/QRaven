package ticketService

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/utils"
	"strconv"
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
		ID:           utils.GenerateTicketId(),
		EventID:      req.EventID,
		UserID:       req.UserID,
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
		PurchaseTime: ticket.PurchaseTime,
		Status:       ticket.Status,
		Amount:       ticket.Amount,
		Type:         ticket.Type,
	}

	return responseData, http.StatusCreated, nil

}

func GetAllTickets(c *gin.Context, db *storage.Database) ([]models.GetTicketResponse, int, error) {

	var tickets []models.Ticket
	var responseData []models.GetTicketResponse
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	offset := (page - 1) * pageSize
	claims, _ := c.Get("userClaims")
	role := claims.(jwt.MapClaims)["role"].(float64)
	id := claims.(jwt.MapClaims)["user_id"].(string)
	fmt.Println(role, "role")
	fmt.Println(id, "id")
	if models.RoleId(role) == models.RoleIdentity.Admin {
		if err := db.Postgresql.Offset(offset).Limit(pageSize).Find(&tickets).Error; err != nil {
			return nil, http.StatusInternalServerError, err
		}
		for _, ticket := range tickets {
			responseData = append(responseData, models.GetTicketResponse{
				ID:           ticket.ID,
				EventID:      ticket.EventID,
				UserID:       ticket.UserID,
				PurchaseTime: ticket.PurchaseTime,
				Status:       ticket.Status,
				Amount:       ticket.Amount,
				Type:         ticket.Type,
			})
		}
		return responseData, http.StatusOK, nil
	}
	if models.RoleId(role) == models.RoleIdentity.User || models.RoleId(role) == models.RoleIdentity.Organizer {
		if err := db.Postgresql.Where("user_id = ?", id).Offset(offset).Limit(pageSize).Find(&tickets).Error; err != nil {
			return nil, http.StatusInternalServerError, err
		}
		for _, ticket := range tickets {
			responseData = append(responseData, models.GetTicketResponse{
				ID:           ticket.ID,
				EventID:      ticket.EventID,
				UserID:       ticket.UserID,
				PurchaseTime: ticket.PurchaseTime,
				Status:       ticket.Status,
				Amount:       ticket.Amount,
				Type:         ticket.Type,
			})

		}
		return responseData, http.StatusOK, nil
	}
	return nil, http.StatusUnauthorized, errors.New("unauthorized")
}
