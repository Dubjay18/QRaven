package router

import (
	"fmt"
	"qraven/pkg/controller/auth"
	"qraven/pkg/repository/storage"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Auth(r *gin.Engine, ApiVersion string, validator *validator.Validate, db *storage.Database, logger *utils.Logger) *gin.Engine {
	auth := auth.Controller{Db: db, Validator: validator, Logger: logger}


	authUrl := r.Group(fmt.Sprintf("%v/auth", ApiVersion))
	{
		authUrl.POST("/register", auth.CreateRegularUser)
		authUrl.POST("/register/admin", auth.CreateAdminUser)
		authUrl.POST("register/organizer", auth.CreateOrganizerUser)
	}



	return r
}
	