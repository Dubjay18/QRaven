package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"qraven/internal/models"
	"qraven/pkg/controller/event"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	"qraven/utils"
)

func Event(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utils.Logger) *gin.Engine {
	event := event.Controller{Db: db, Validator: validator, Logger: logger}

	eventUrl := r.Group(fmt.Sprintf("%v/events", ApiVersion))
	{
		eventUrl.POST("/", middleware.Authorize(db.Postgresql, models.RoleIdentity.Organizer), event.CreateEvent)
		eventUrl.GET("/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.Organizer, models.RoleIdentity.User, models.RoleIdentity.Admin), event.GetEvent)
		eventUrl.GET("/", middleware.Authorize(db.Postgresql, models.RoleIdentity.Organizer, models.RoleIdentity.User, models.RoleIdentity.Admin), event.GetAllEvents)
		eventUrl.PUT("/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.Organizer), event.UpdateEvent)
	}

	return r
}
