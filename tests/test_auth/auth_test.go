package testauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"qraven/internal/models"
	"qraven/pkg/controller/auth"
	"qraven/pkg/repository/storage"
	"qraven/tests"
	"qraven/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TestUserSignup(t *testing.T) {
	logger := tests.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	requestURI := url.URL{Path: tests.SignupPath}
	currUUID := utils.GenerateUUID()

	signupTests := []struct {
		Name         string
		RequestBody  models.CreateUserRequest
		ExpectedCode int
		Message      string
	}{
		{
			Name: "Successful user register",
			RequestBody: models.CreateUserRequest{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				Gender:      "male",
				DateOfBirth: "1990-01-01",
				Role:        "user",
			},
			ExpectedCode: http.StatusCreated,
			Message:      "user created successfully",
		}, {
			Name: "details already exist",
			RequestBody: models.CreateUserRequest{
				Email:       fmt.Sprintf("testuser%v@qa.team", currUUID),
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				Gender:      "male",
				DateOfBirth: "1990-01-01",
				Role:        "user",
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "user already exists with the given email",
		}, {
			Name: "invalid email",
			RequestBody: models.CreateUserRequest{
				Email:       "emailtest",
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				Gender:      "male",
				DateOfBirth: "1990-01-01",
				Role:        "user",
			},
			ExpectedCode: http.StatusBadRequest,
			Message:      "email address is invalid",
		}, {
			Name: "Validation failed",
			RequestBody: models.CreateUserRequest{
				FirstName:   "test",
				LastName:    "user",
				Password:    "password",
				Gender:      "male",
				DateOfBirth: "1990-01-01",
				Role:        "male",
			},
			ExpectedCode: http.StatusUnprocessableEntity,
			Message:      "Validation failed",
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}

	for _, test := range signupTests {
		r := gin.Default()

		r.POST(tests.SignupPath, auth.CreateUser)

		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, requestURI.String(), &b)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			tests.AssertStatusCode(t, rr.Code, test.ExpectedCode)

			data := tests.ParseResponse(rr)

			code := int(data["status_code"].(float64))
			tests.AssertStatusCode(t, code, test.ExpectedCode)

			if test.Message != "" {
				message := data["message"]
				if message != nil {
					tests.AssertResponseMessage(t, message.(string), test.Message)
				} else {
					tests.AssertResponseMessage(t, "", test.Message)
				}

			}

		})

	}

}
