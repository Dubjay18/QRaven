package event

import (
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	eventService "qraven/services/auth/event"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)


type Controller struct {
	Db 	  *storage.Database
	Validator *validator.Validate
	Logger *utils.Logger
}

func (base *Controller) CreateEvent(c *gin.Context) {
	// create event
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var rd utils.Response

		if ve, ok := err.(validator.ValidationErrors); ok {
			rd = utils.BuildErrorResponse(http.StatusBadRequest, "error","failed to validate request",utils.ValidationResponse(ve, base.Validator,req),nil)

			c.JSON(http.StatusBadRequest, rd)
			return
		}
		rd = utils.BuildErrorResponse(http.StatusBadRequest, "error","failed to parse request",err,nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// validate request
	err := base.Validator.Struct(req);
	if  err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "eror","failed to validate request",utils.ValidationResponse(err, base.Validator,req),nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}
	// err = eventService.ValidateRequest(req, base.Db.Postgresql);
	// if  err != nil {
		// rd := utils.BuildErrorResponse(http.StatusBadRequest, "error","Bad request",err,nil)
		// c.JSON(http.StatusBadRequest, rd)
		// return
	// }

	if res, code, err := eventService.CreateEvent(req, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error","failed to create event",err,nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "event created successfully", res)
		c.JSON(code, rd)
	}
}
