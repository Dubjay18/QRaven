package router

import (
	"fmt"
	"qraven/internal/models"
	"qraven/pkg/controller/payment"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Payment(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utils.Logger) *gin.Engine {
	paymentController := payment.Controller{Db: db, Validator: validator, Logger: logger}

	paymentURL := r.Group(fmt.Sprintf("%v/payments", ApiVersion))
	{
		paymentURL.POST("/initialize", middleware.Authorize(db.Postgresql, models.RoleIdentity.User, models.RoleIdentity.Organizer, models.RoleIdentity.Admin), paymentController.InitializePayment)
		paymentURL.GET("/:id", middleware.Authorize(db.Postgresql, models.RoleIdentity.User, models.RoleIdentity.Organizer, models.RoleIdentity.Admin), paymentController.GetPaymentByID)
		paymentURL.PATCH("/:id/status", middleware.Authorize(db.Postgresql, models.RoleIdentity.User, models.RoleIdentity.Organizer, models.RoleIdentity.Admin), paymentController.UpdatePaymentStatus)
	}

	return r
}
