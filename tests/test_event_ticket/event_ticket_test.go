package testeventticket

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"qraven/internal/models"
	"qraven/pkg/repository/storage"
	"qraven/pkg/router"
	"qraven/tests"
	"qraven/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func setupEventTicketRouter(t *testing.T) (*gin.Engine, *storage.Database) {
	if os.Getenv("DB_HOST") == "" || os.Getenv("DB_PORT") == "" || os.Getenv("USERNAME") == "" || os.Getenv("DB_NAME") == "" {
		t.Skip("Skipping event/ticket integration tests: DB environment variables are not configured")
	}

	logger := tests.Setup()
	gin.SetMode(gin.TestMode)
	validatorRef := validator.New()
	db := storage.Connection()

	r := gin.Default()
	apiVersion := "api/v1"
	router.Auth(r, apiVersion, validatorRef, db, logger)
	router.Event(r, apiVersion, validatorRef, db, logger)
	router.Ticket(r, apiVersion, validatorRef, db, logger)

	return r, db
}

func registerAndLogin(t *testing.T, r *gin.Engine, role string) (string, string) {
	currUUID := utils.GenerateUUID()
	email := "acct-" + currUUID + "@qa.team"

	signupBody := models.CreateUserRequest{
		Email:       email,
		FirstName:   "test",
		LastName:    "user",
		Password:    "password",
		Gender:      "male",
		DateOfBirth: "1990-01-01",
		Role:        role,
	}

	var signupBuffer bytes.Buffer
	json.NewEncoder(&signupBuffer).Encode(signupBody)
	signupReq, err := http.NewRequest(http.MethodPost, "/api/v1/auth/register", &signupBuffer)
	if err != nil {
		t.Fatal(err)
	}
	signupReq.Header.Set("Content-Type", "application/json")

	signupResp := httptest.NewRecorder()
	r.ServeHTTP(signupResp, signupReq)
	if signupResp.Code != http.StatusCreated {
		t.Fatalf("signup failed: expected %d got %d body=%s", http.StatusCreated, signupResp.Code, signupResp.Body.String())
	}

	signupData := tests.ParseResponse(signupResp)
	dataObj, ok := signupData["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("invalid signup response data")
	}
	userObj, ok := dataObj["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("invalid signup user object")
	}
	userID, ok := userObj["id"].(string)
	if !ok || userID == "" {
		t.Fatalf("unable to resolve user id from signup response")
	}

	loginBody := models.UserLoginRequest{Email: email, Password: "password"}
	var loginBuffer bytes.Buffer
	json.NewEncoder(&loginBuffer).Encode(loginBody)

	loginReq, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &loginBuffer)
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp := httptest.NewRecorder()
	r.ServeHTTP(loginResp, loginReq)
	if loginResp.Code != http.StatusOK {
		t.Fatalf("login failed: expected %d got %d body=%s", http.StatusOK, loginResp.Code, loginResp.Body.String())
	}

	loginData := tests.ParseResponse(loginResp)
	loginDataObj, ok := loginData["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("invalid login response data")
	}
	accessToken, ok := loginDataObj["access_token"].(string)
	if !ok || accessToken == "" {
		t.Fatalf("unable to resolve access token from login response")
	}

	return accessToken, userID
}

func createEventAsOrganizer(t *testing.T, r *gin.Engine, organizerToken, organizerID string) string {
	body := models.CreateEventRequest{
		Title:       "Event " + utils.GenerateUUID(),
		Description: "integration event",
		StartDate:   "2026-08-01",
		EndDate:     "2026-08-02",
		Location:    "Lagos",
		TicketPrice: 100,
		Capacity:    100,
		OrganizerID: organizerID,
	}

	var requestBody bytes.Buffer
	json.NewEncoder(&requestBody).Encode(body)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/events/", &requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+organizerToken)

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("create event failed: expected %d got %d body=%s", http.StatusCreated, resp.Code, resp.Body.String())
	}

	response := tests.ParseResponse(resp)
	dataObj, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("invalid create event response data")
	}
	eventID, ok := dataObj["id"].(string)
	if !ok || eventID == "" {
		t.Fatalf("unable to resolve event id from create event response")
	}

	return eventID
}

func TestOrganizerCanCreateAndListEvents(t *testing.T) {
	r, _ := setupEventTicketRouter(t)
	organizerToken, organizerID := registerAndLogin(t, r, string(models.OrganizerRole))
	_ = createEventAsOrganizer(t, r, organizerToken, organizerID)

	listReq, err := http.NewRequest(http.MethodGet, (&url.URL{Path: "/api/v1/events/"}).String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	listReq.Header.Set("Authorization", "Bearer "+organizerToken)

	listResp := httptest.NewRecorder()
	r.ServeHTTP(listResp, listReq)

	tests.AssertStatusCode(t, listResp.Code, http.StatusOK)
	listData := tests.ParseResponse(listResp)
	statusCode := int(listData["status_code"].(float64))
	tests.AssertStatusCode(t, statusCode, http.StatusOK)
}

func TestUserCanCreateAndListOwnTickets(t *testing.T) {
	r, _ := setupEventTicketRouter(t)
	organizerToken, organizerID := registerAndLogin(t, r, string(models.OrganizerRole))
	userToken, userID := registerAndLogin(t, r, string(models.UserRole))
	eventID := createEventAsOrganizer(t, r, organizerToken, organizerID)

	createTicketBody := models.CreateTicketRequest{
		EventID: eventID,
		UserID:  userID,
		Amount:  1,
		Type:    models.TicketTypeRegular,
	}

	var createTicketBuffer bytes.Buffer
	json.NewEncoder(&createTicketBuffer).Encode(createTicketBody)

	createTicketReq, err := http.NewRequest(http.MethodPost, "/api/v1/tickets/"+eventID, &createTicketBuffer)
	if err != nil {
		t.Fatal(err)
	}
	createTicketReq.Header.Set("Content-Type", "application/json")
	createTicketReq.Header.Set("Authorization", "Bearer "+userToken)

	createTicketResp := httptest.NewRecorder()
	r.ServeHTTP(createTicketResp, createTicketReq)
	tests.AssertStatusCode(t, createTicketResp.Code, http.StatusCreated)

	listTicketReq, err := http.NewRequest(http.MethodGet, "/api/v1/tickets/", nil)
	if err != nil {
		t.Fatal(err)
	}
	listTicketReq.Header.Set("Authorization", "Bearer "+userToken)

	listTicketResp := httptest.NewRecorder()
	r.ServeHTTP(listTicketResp, listTicketReq)
	tests.AssertStatusCode(t, listTicketResp.Code, http.StatusOK)

	listTicketData := tests.ParseResponse(listTicketResp)
	statusCode := int(listTicketData["status_code"].(float64))
	tests.AssertStatusCode(t, statusCode, http.StatusOK)
}
