package authService

import (
	"errors"
	"fmt"
	"net/http"
	"qraven/internal/models"
	"qraven/pkg/middleware"
	"qraven/pkg/repository/storage/postgresql"
	"qraven/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func ValidateRequest(req models.CreateUserRequest, db *gorm.DB) (error) {
if ok := postgresql.CheckExistsInTable(db, "users", "email = ?", req.Email); ok {
		return  errors.New("user with this email already exists")
	}
	return  nil
}

func CreateUser(req models.CreateUserRequest, role models.RoleId, db *gorm.DB) (gin.H, int, error) {
	user := models.User{}
	var responseData gin.H
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, http.StatusInternalServerError,err
	}
	user = models.User{
		ID: utils.GenerateUUID(),
		FirstName: strings.ToLower(req.FirstName),
		LastName: strings.ToLower(req.LastName),
		Email: strings.ToLower(req.Email),
		Password:hashedPassword,
		Role: role,
	}
err = user.CreateUser(db)
if err != nil {
	return nil, http.StatusInternalServerError,err
}

tokenData, err := middleware.CreateToken(user)
	if err != nil {
	
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}

	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: user.ID}

	err = access_token.CreateAccessToken(db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}
	reponseData := gin.H{
		"user": models.CreateUserResponse{
			ID: user.ID,
			FirstName: user.FirstName,
			LastName: user.LastName,
			Email: user.Email,
			Role: string(user.GetRoleName()),
		},
		"access_token": tokenData.AccessToken,
	}

	return reponseData, http.StatusCreated, nil  
}