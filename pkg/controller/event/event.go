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
			return
		}
		rd = utils.BuildErrorResponse(http.StatusBadRequest, "error", "failed to parse request", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// validate request
	err := base.Validator.Struct(req)
	if err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "eror", "failed to validate request", utils.ValidationResponse(err, base.Validator, req), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	if res, code, err := eventService.CreateEvent(req, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to create event", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "event created successfully", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) GetEvent(c *gin.Context) {
	// get event
	eventId := c.Params.ByName("id")

	if res, code, err := eventService.GetEventByID(eventId, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to get event", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "event retrieved successfully", res)
		c.JSON(code, rd)
	}
}