package payment

import (
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	paymentService "qraven/services/payment"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utils.Logger
}

func (base *Controller) InitializePayment(c *gin.Context) {
	var req models.InitializePaymentRequest
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

	err := base.Validator.Struct(req)
	if err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "failed to validate request", utils.ValidationResponse(err, base.Validator, req), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	if res, code, err := paymentService.InitializePayment(req, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to initialize payment", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "payment initialized successfully", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) GetPaymentByID(c *gin.Context) {
	paymentID := c.Param("id")

	if res, code, err := paymentService.GetPaymentByID(paymentID, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to get payment", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "payment retrieved successfully", res)
		c.JSON(code, rd)
	}
}

func (base *Controller) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")
	var req models.UpdatePaymentStatusRequest

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

	err := base.Validator.Struct(req)
	if err != nil {
		rd := utils.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "failed to validate request", utils.ValidationResponse(err, base.Validator, req), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	if res, code, err := paymentService.UpdatePaymentStatus(paymentID, req, base.Db); err != nil {
		rd := utils.BuildErrorResponse(code, "error", "failed to update payment status", err, nil)
		c.JSON(code, rd)
		return
	} else {
		rd := utils.BuildSuccessResponse(code, "payment status updated successfully", res)
		c.JSON(code, rd)
	}
}
