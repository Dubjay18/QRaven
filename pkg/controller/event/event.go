package event

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	eventService "qraven/services/auth/event"
	"qraven/utils"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utils.Logger
}

func (base *Controller) CreateEvent(c *gin.Context) {
	// create event
	var req models.CreateEventRequest
	req.OrganizerID, _ = middleware.GetIdFromToken(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		var rd utils.Response

		if ve, ok := err.(validator.ValidationErrors); ok {
			rd = utils.BuildErrorResponse(http.StatusBadRequest, "error", "failed to validate request", utils.ValidationResponse(ve, base.Validator, req), nil)
			c.JSON(http.StatusBadRequest, rd)
			base.Logger.Error(err)
			return
		}
		rd = utils.BuildErrorResponse(http.StatusBadRequest, "error", "failed to parse request", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// validate request
	err := base.Validator.Struct(req)
	if err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "failed to validate request", utils.ValidationResponse(err, base.Validator, req), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		base.Logger.Error(err)
		return
	}

	if res, code, err := eventService.CreateEvent(req, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to create event", err, nil)
		base.Logger.Error(err)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "event created successfully", res)
		base.Logger.Info("event retrieved successfully")
		c.JSON(code, rd)
	}
}

func (base *Controller) GetEvent(c *gin.Context) {
	// get event
	eventId := c.Params.ByName("id")
	if res, code, err := eventService.GetEventByID(eventId, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to get event", err, nil)
		base.Logger.Error(err)
		c.JSON(code, rd)
		return
	} else {
		base.Logger.Info("event retrieved successfully")
		rd := utils.BuildSuccessResponse(code, "event retrieved successfully", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) GetAllEvents(c *gin.Context) {
	if res, code, err := eventService.GetEvents(base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to get events", err, nil)
		base.Logger.Error(err)
		c.JSON(code, rd)
		return
	} else {
		base.Logger.Info("events retrieved successfully")
		rd := utils.BuildSuccessResponse(code, "events retrieved successfully", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) UpdateEvent(c *gin.Context) {
	// update event
	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var rd utils.Response

		if ve, ok := err.(validator.ValidationErrors); ok {
			rd = utils.BuildErrorResponse(http.StatusBadRequest, "error", "failed to validate request", utils.ValidationResponse(ve, base.Validator, req), nil)
			c.JSON(http.StatusBadRequest, rd)
			base.Logger.Error(err)
			return
		}
		rd = utils.BuildErrorResponse(http.StatusBadRequest, "error", "failed to parse request", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// validate request
	err := base.Validator.Struct(req)
	if err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "failed to validate request", utils.ValidationResponse(err, base.Validator, req), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		base.Logger.Error(err)
		return
	}

	eventId := c.Params.ByName("id")
	if res, code, err := eventService.UpdateEventData(req, eventId, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to update event", err, nil)
		base.Logger.Error(err)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "event updated successfully", res)
		base.Logger.Info("event updated successfully")
		c.JSON(code, rd)
	}

}

func (base *Controller) DeleteEvent(c *gin.Context) {
	// delete event
	eventId := c.Params.ByName("id")
	if code, err := eventService.DeleteEvent(eventId, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to delete event", err, nil)
		base.Logger.Error(err)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "event deleted successfully", nil)
		base.Logger.Info("event deleted successfully")
		c.JSON(code, rd)
	}
}
