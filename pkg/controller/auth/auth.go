package auth

import (
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	authService "qraven/services/auth"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utils.Logger
}

func (base *Controller) CreateUser(c *gin.Context) {
	// create user
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rd := utils.BuildErrorResponse(http.StatusBadRequest, "error","failed to parse request",err,nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// validate request
	if err := base.Validator.Struct(req); err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "error","failed to validate request",err,nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	if err := authService.ValidateRequest(req, base.Db.Postgresql); err != nil {
		rd := utils.BuildErrorResponse(http.StatusBadRequest, "error","failed to validate request",err,nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if res, code, err := authService.CreateUser(req, base.Db.Postgresql); err != nil {
		rd := utils.BuildErrorResponse(code, "error","failed to create user",err,nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "user created successfully", res)
		c.JSON(code, rd)
	}
}