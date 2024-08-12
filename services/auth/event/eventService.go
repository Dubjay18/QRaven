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
	if ok := postgresql.CheckExistsInTable(db.Postgresql, "events", "title = ?, location = ?, organizer_id = ?", req.Title, req.Location, req.OrganizerID); ok {
		return models.CreateEventResponse{}, http.StatusConflict, errors.New("event with this title, location and organizer already exists")
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
