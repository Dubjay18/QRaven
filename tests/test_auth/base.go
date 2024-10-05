package testauth

// import (
// 	"qraven/pkg/controller/auth"
// 	"qraven/pkg/repository/storage"
// 	"qraven/tests"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// )

// func SetupAuthTestRouter(*gin.Engine, *auth.Controller) {
// 	gin.SetMode(gin.TestMode)

// 	logger := tests.Setup()
// 	db := storage.Connection()
// 	validator := validator.New()

// 	authController := &auth.Controller{
// 		Db:        db,
// 		Validator: validator,
// 		Logger:    logger,
// 	}

// 	r := gin.Default()
// 	SetupAuthRoutes(r, authController)
// 	return r, authController
// }

// func SetupAuthRoutes(r *gin.Engine, authController *auth.Controller) {

// 	r.POST("/api/v1/auth/password-reset", authController.ResetPassword)
// 	r.POST("/api/v1/auth/password-reset/verify", authController.VerifyResetToken)
// 	r.POST("/api/v1/auth/magick-link", authController.RequestMagicLink)
// 	r.POST("/api/v1/auth/magick-link/verify", authController.VerifyMagicLink)
// }
