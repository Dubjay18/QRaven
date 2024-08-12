package router

import (
	"fmt"
	"qraven/internal/models"
	"qraven/pkg/controller/event"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Event(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utils.Logger) *gin.Engine {
	event := event.Controller{Db: db, Validator: validator, Logger: logger}

	eventUrl := r.Group(fmt.Sprintf("%v/event", ApiVersion))
	{
		eventUrl.POST("/create", middleware.Authorize(db.Postgresql, models.RoleIdentity.Organizer), event.CreateEvent)
	}

	return r
}
