package eventService

import (
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"qraven/utils"
)



func CreateEvent(req models.CreateEventRequest, db *storage.Database) (models.CreateEventResponse, int, error) {
	// create event
	event := models.Event{}
	var responseData models.CreateEventResponse
	event = models.Event{
		ID:          utils.GenerateUUID(),
		Title:        req.Title,
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