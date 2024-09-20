package eventService

import (
	"errors"
	"fmt"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context, req models.CreateEventRequest, db *storage.Database) (models.CreateEventResponse, int, error) {
	//check if event exists
	if ok := postgresql.CheckExistsInTable(db.Postgresql, "events", "title = ? AND location = ? AND organizer_id = ?", req.Title, req.Location, req.OrganizerID); ok {
		return models.CreateEventResponse{}, http.StatusConflict, errors.New("event with this title, location, and organizer already exists")
	}

	imageFile, _ := c.FormFile("image")
	if imageFile != nil {
		// Upload the image file
		imagePath, err := utils.UploadFile(imageFile)
		if err != nil {
			return models.CreateEventResponse{}, http.StatusInternalServerError, err
		}
		req.Image = imagePath
	} else {
		req.Image = "../../static/logo.webp"
	}

	// create event
	event := models.Event{}
	var responseData models.CreateEventResponse
	event = models.Event{
		ID:          utils.GenerateUUID(),
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    req.Location,
		TicketPrice: req.TicketPrice,
		Capacity:    req.Capacity,
		OrganizerID: req.OrganizerID,
	}
	err := event.CreateEvent(db.Postgresql)
	if err != nil {
		return responseData, http.StatusInternalServerError, err
	}

	responseData = models.CreateEventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		TicketPrice: event.TicketPrice,
		Capacity:    event.Capacity,
		OrganizerID: event.OrganizerID,
	}

	return responseData, http.StatusCreated, nil

}

func GetEventByID(eventId string, db *storage.Database) (models.GetEventResponse, int, error) {
	var responseData models.GetEventResponse

	event := models.Event{
		ID: eventId,
	}
	err := event.GetEventByID(db.Postgresql)
	if err != nil {
		return responseData, http.StatusNotFound, err
	}

	return models.GetEventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		TicketPrice: event.TicketPrice,
		Capacity:    event.Capacity,
		OrganizerID: event.OrganizerID,
	}, http.StatusOK, nil

}

func GetEvents(c *gin.Context, db *storage.Database) ([]models.GetEventResponse, int, error) {
	var responseData []models.GetEventResponse
	var events []models.Event
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	offset := (page - 1) * pageSize
	if err := db.Postgresql.Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, http.StatusInternalServerError, err
	}

	fmt.Println(events, "fjhfj")
	for _, event := range events {
		responseData = append(responseData, models.GetEventResponse{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    event.Location,
			TicketPrice: event.TicketPrice,
			Capacity:    event.Capacity,
			OrganizerID: event.OrganizerID,
			TicketCount: event.TicketCount,
		})
	}

	return responseData, http.StatusOK, nil

}

func UpdateEventData(updateData models.UpdateEventRequest, eventId string, db *storage.Database) (models.GetEventResponse, int, error) {
	var responseData models.GetEventResponse
	event := models.Event{
		ID: eventId,
	}
	err := event.GetEventByID(db.Postgresql)
	if err != nil {
		return responseData, http.StatusNotFound, err
	}
	if updateData.Title != "" {
		event.Title = updateData.Title
	}
	if updateData.Description != "" {
		event.Description = updateData.Description
	}
	if updateData.StartDate != "" {
		event.StartDate = updateData.StartDate
	}
	if updateData.EndDate != "" {
		event.EndDate = updateData.EndDate
	}
	if updateData.Location != "" {
		event.Location = updateData.Location
	}
	if updateData.TicketPrice != 0 {
		event.TicketPrice = updateData.TicketPrice
	}
	if updateData.Capacity != 0 {
		event.Capacity = updateData.Capacity
	}
	if updateData.OrganizerID != "" {
		event.OrganizerID = updateData.OrganizerID
	}
	err = event.UpdateEvent(db.Postgresql)
	if err != nil {
		return responseData, http.StatusInternalServerError, err
	}

	return models.GetEventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		TicketPrice: event.TicketPrice,
		Capacity:    event.Capacity,
		OrganizerID: event.OrganizerID,
	}, http.StatusOK, nil

}

func DeleteEvent(eventId string, db *storage.Database) (int, error) {
	event := models.Event{
		ID: eventId,
	}
	err := event.GetEventByID(db.Postgresql)
	if err != nil {
		return http.StatusNotFound, err
	}
	err = event.DeleteEvent(db.Postgresql)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}
