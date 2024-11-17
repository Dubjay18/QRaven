package testauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"net/url"
	"qraven/internal/models"
	"qraven/pkg/controller/auth"
	"qraven/pkg/repository/storage"
	"qraven/tests"
	"qraven/utils"
	"testing"
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

func TestUserLogin(t *testing.T) {
	logger := tests.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()
	var (
		loginPath = "/api/v1/auth/login"
		loginURI  = url.URL{Path: loginPath}
		// currUUID       = utils.GenerateUUID()
		userSignUpData = models.CreateUserRequest{
			Email:       "testuser@qa.team",
			FirstName:   "test",
			LastName:    "user",
			Password:    "password",
			Gender:      "female",
			DateOfBirth: "1990-01-01",
			Role:        "user",
		}
	)

	loginTests := []struct {
		Name         string
		RequestBody  models.UserLoginRequest
		ExpectedCode int
		Message      string
	}{
		{
			Name: "OK email login successful",
			RequestBody: models.UserLoginRequest{
				Email:    userSignUpData.Email,
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusOK,
			Message:      "user login successfully",
		}, {
			Name:         "password not provided",
			RequestBody:  models.UserLoginRequest{},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "username or phone or email not provided",
			RequestBody: models.UserLoginRequest{
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "email does not exist",
			RequestBody: models.UserLoginRequest{
				Email:    utils.GenerateUUID(),
				Password: userSignUpData.Password,
			},
			ExpectedCode: http.StatusBadRequest,
		}, {
			Name: "incorrect password",
			RequestBody: models.UserLoginRequest{
				Email:    "testuser@qa.team",
				Password: "incorrect",
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	auth := auth.Controller{Db: db, Validator: validatorRef, Logger: logger}
	r := gin.Default()
	tests.SignupUser(t, r, auth, userSignUpData)
	r.POST(loginPath, auth.Login)

	for _, test := range loginTests {
		t.Run(test.Name, func(t *testing.T) {
			var b bytes.Buffer
			json.NewEncoder(&b).Encode(test.RequestBody)

			req, err := http.NewRequest(http.MethodPost, loginURI.String(), &b)
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
