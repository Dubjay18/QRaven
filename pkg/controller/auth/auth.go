package auth

import (
	"log"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	authService "qraven/services/auth"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utils.Logger
}

func (base *Controller) CreateRegularUser(c *gin.Context) {
	// create user
	var req models.CreateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		var rd utils.Response

		if ve, ok := err.(validator.ValidationErrors); ok {
			rd = utils.BuildErrorResponse(http.StatusBadRequest, "error", "failed to validate request", utils.ValidationResponse(ve, base.Validator, req), nil)

			c.JSON(http.StatusBadRequest, rd)
			return
		}
		// base.Logger.Error(err)
		log.Println(err)
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
	err = authService.ValidateRequest(req, base.Db.Postgresql)
	if err != nil {
		rd := utils.BuildErrorResponse(http.StatusBadRequest, "error", "Bad request", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if res, code, err := authService.CreateUser(c, req, base.Db.Postgresql); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to create user", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "user created successfully", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) Login(c *gin.Context) {
	var req models.UserLoginRequest
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
	if err := base.Validator.Struct(req); err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "failed to validate request", utils.ValidationResponse(err, base.Validator, req), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	if res, code, err := authService.Login(req, base.Db.Postgresql); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to login", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "login successful", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) LogoutUser(c *gin.Context) {
	claims, exists := c.Get("userClaims")
	if !exists {
		rd := utils.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get user claims", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userClaims := claims.(jwt.MapClaims)

	access_uuid, _ := userClaims["access_uuid"].(string)
	owner_id, ok := userClaims["user_id"].(string)
	if !ok {
		rd := utils.BuildErrorResponse(http.StatusBadRequest, "error", "unable to get access id", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	respData, code, err := authService.LogoutUser(access_uuid, owner_id, base.Db.Postgresql)
	if err != nil {
		rd := utils.BuildErrorResponse(code, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("user logout successfully")

	rd := utils.BuildSuccessResponse(http.StatusOK, "user logout successfully", respData)
	c.JSON(http.StatusOK, rd)
}
