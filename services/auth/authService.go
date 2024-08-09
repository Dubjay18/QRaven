package authService

import (
	"errors"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func ValidateRequest(req models.CreateUserRequest, db *gorm.DB) (error) {
if ok := postgresql.CheckExistsInTable(db, "users", "email = ?", req.Email); ok {
		return  errors.New("failed to check user existence")
	}
	return  nil
}

func CreateUser(req models.CreateUserRequest, db *gorm.DB) (gin.H, int, error) {
	user := models.User{}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, http.StatusInternalServerError,err
	}
	user = models.User{
		ID: utils.GenerateUUID(),
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email: req.Email,
		Password:hashedPassword,
		Role: models.UserRole,
	}
err = user.CreateUser(db)
if err != nil {
	return nil, http.StatusInternalServerError,err
}
	reponseData := gin.H{
		"user": models.CreateUserResponse{
			ID: user.ID,
			FirstName: user.FirstName,
			LastName: user.LastName,
			Email: user.Email,
			Role: user.Role,
		},
	}

	return reponseData, http.StatusCreated, nil  
}