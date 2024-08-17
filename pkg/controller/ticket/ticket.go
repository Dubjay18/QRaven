package ticket

import (
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	ticketService "qraven/services/ticket"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utils.Logger
}

// CreateTicket creates a new ticket
// @Summary      Creates a new ticket for an event
// @Description  Creates a new ticket for the event with the given eventId.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        eventId   path      string                     true  "Event ID"
// @Param        request   body      models.CreateTicketRequest true  "Create Ticket Request"
// @Success      200       {object}  models.CreateTicketResponse  "ticketResponse"
// @Failure      400       {object}  utils.Response              "badRequest"
// @Failure      422       {object}  utils.Response              "validationError"
// @Router       /ticket/{eventId} [post]
func (base *Controller) CreateTicket(c *gin.Context) {
	// create ticket
	var req models.CreateTicketRequest
	req.EventID = c.Param("eventId")
	req.UserID, _ = middleware.GetIdFromToken(c)
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

	if res, code, err := ticketService.CreateTicket(req, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to create ticket", err, nil)
		// base.Logger.Error(err)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "ticket created successfully", res)
		base.Logger.Info("ticket created successfully")
		c.JSON(code, rd)
	}
}

func (base *Controller) GetTickets(c *gin.Context) {
	if res, code, err := ticketService.GetAllTickets(c, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to get tickets", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "tickets retrieved successfully", res)
		c.JSON(code, rd)
	}

}
