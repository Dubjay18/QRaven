package router

import (
	"fmt"
	"qraven/internal/models"
	"qraven/pkg/controller/ticket"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Ticket(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utils.Logger) *gin.Engine {

	ticket := ticket.Controller{Db: db, Validator: validator, Logger: logger}

	ticketUrl := r.Group(fmt.Sprintf("%v/tickets", ApiVersion))
	{
		ticketUrl.POST("/:eventId", middleware.Authorize(db.Postgresql, models.RoleIdentity.User, models.RoleIdentity.Organizer, models.RoleIdentity.Admin), ticket.CreateTicket)
		// ticketUrl.GET("/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.User, models.RoleIdentity.Admin), ticket.GetTicket)
		// ticketUrl.GET("/", middleware.Authorize(db.Postgresql, models.RoleIdentity.User, models.RoleIdentity.Admin), ticket.GetAllTickets)
		// ticketUrl.PUT("/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.User), ticket.UpdateTicket)
		// ticketUrl.DELETE("/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.User), ticket.DeleteTicket)
	}

	return r
}
