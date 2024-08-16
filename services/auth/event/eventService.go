package eventService

import (
	"errors"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/utils"
)

func CreateEvent(req models.CreateEventRequest, db *storage.Database) (models.CreateEventResponse, int, error) {
	//check if event exists
	if ok := postgresql.CheckExistsInTable(db.Postgresql, "events", "title = ? AND location = ? AND organizer_id = ?", req.Title, req.Location, req.OrganizerID); ok {
		return models.CreateEventResponse{}, http.StatusConflict, errors.New("event with this title, location, and organizer already exists")
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

func GetEvents(db *storage.Database) ([]models.GetEventResponse, int, error) {
	var responseData []models.GetEventResponse
	err := postgresql.SelectAllFromDb(db.Postgresql, "desc", &responseData, nil)
	if err != nil {
		return responseData, http.StatusInternalServerError, err
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
