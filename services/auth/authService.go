package authService

import (
	"errors"
	"fmt"
	"log"
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

func ValidateRequest(req models.CreateUserRequest, db *gorm.DB) error {
	if ok := postgresql.CheckExistsInTable(db, "users", "email = ?", req.Email); ok {
		return errors.New("user already exists with the given email")
	}
	if req.Gender != "" && req.Gender != "male" && req.Gender != "female" {
		return errors.New("invalid gender")
	}
	if !utils.ValidateEmail(req.Email) {
		return errors.New("email address is invalid")
	}
	if req.DateOfBirth == "" {
		return errors.New("date of birth is required")
	}
	return nil
}

func CreateUser(c *gin.Context, req models.CreateUserRequest, db *gorm.DB) (gin.H, int, error) {
	// Create a new user
	user := models.User{}
	var responseData gin.H
	if req.Role == "" || req.Role != string(models.UserRole) && req.Role != string(models.OrganizerRole) {
		return responseData, http.StatusBadRequest, errors.New("invalid role")
	}
	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Parse the date of birth
	parsedDateOfBirth, err := req.ParseDateOfBirth()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Handle file upload for the avatar
	avatarFile, _ := c.FormFile("avatar")
	if avatarFile != nil {
		// Upload the avatar file
		avatarPath, err := utils.UploadFile(avatarFile)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		log.Println(avatarPath)
		req.Avatar = avatarPath
	} else {
		// Generate a default avatar URL
		req.Avatar = "https://www.gravatar.com/avatar/" + utils.GenerateUUID()
	}

	// Set the user data
	user = models.User{
		ID:          utils.GenerateUUID(),
		FirstName:   strings.ToLower(req.FirstName),
		LastName:    strings.ToLower(req.LastName),
		Email:       strings.ToLower(req.Email),
		Password:    hashedPassword,
		Role:        models.RoleName(req.Role),
		Gender:      req.Gender,
		DateOfBirth: parsedDateOfBirth,
		Avatar:      req.Avatar,
	}

	// Create the user in the database
	err = user.CreateUser(db)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Create an access token for the user
	tokenData, err := middleware.CreateToken(user)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	// Store the access token in the database
	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}
	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: user.ID}
	err = access_token.CreateAccessToken(db, tokens)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	// Prepare the response data
	reponseData := gin.H{
		"user": models.CreateUserResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Role:        string(user.GetRoleName()),
			Gender:      user.Gender,
			Avatar:      user.Avatar,
			DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
		},
		"access_token": tokenData.AccessToken,
	}

	return reponseData, http.StatusCreated, nil
}

func Login(req models.UserLoginRequest, db *gorm.DB) (gin.H, int, error) {
	var (
		user         = models.User{}
		responseData gin.H
	)
	// Check if the user email exists

	if ok := postgresql.CheckExistsInTable(db, "users", "email = ?", req.Email); !ok {
		return responseData, http.StatusBadRequest, fmt.Errorf("invalid credentials")
	}

	userData, err := user.GetUserByEmail(db, req.Email)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}
	log.Println(userData.Password, user.Password)

	if !utils.CompareHash(req.Password, userData.Password) {
		return nil, http.StatusBadRequest, errors.New("invalid password")
	}
	tokenData, err := middleware.CreateToken(*userData)
	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}

	tokens := map[string]string{
		"access_token": tokenData.AccessToken,
		"exp":          strconv.Itoa(int(tokenData.ExpiresAt.Unix())),
	}

	access_token := models.AccessToken{ID: tokenData.AccessUuid, OwnerID: userData.ID}

	err = access_token.CreateAccessToken(db, tokens)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error saving token: " + err.Error())
	}
	reponseData := gin.H{
		"user": models.CreateUserResponse{
			ID:        userData.ID,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Email:     userData.Email,
			Role:      string(userData.GetRoleName()),
		},
		"access_token": tokenData.AccessToken,
	}

	return reponseData, http.StatusOK, nil
}

func LogoutUser(access_uuid, owner_id string, db *gorm.DB) (gin.H, int, error) {

	var (
		responseData gin.H
	)

	access_token := models.AccessToken{ID: access_uuid, OwnerID: owner_id}

	// revoke user access_token to invalidate session
	err := access_token.RevokeAccessToken(db)

	if err != nil {
		return responseData, http.StatusInternalServerError, fmt.Errorf("error revoking user session: " + err.Error())
	}

	responseData = gin.H{}

	return responseData, http.StatusOK, nil
}
