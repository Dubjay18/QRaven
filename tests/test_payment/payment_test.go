package testpayment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"qraven/internal/models"
	paymentController "qraven/pkg/controller/payment"
	"qraven/pkg/repository/storage"
	"qraven/tests"
	"qraven/utils"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func seedTicketAndPayment(t *testing.T, db *storage.Database) (string, string) {
	userID := utils.GenerateUUID()
	eventID := utils.GenerateUUID()
	ticketID := utils.GenerateTicketId()
	paymentID := "pay_" + utils.GenerateUUID()

	user := models.User{
		ID:          userID,
		FirstName:   "Payment",
		LastName:    "Tester",
		Email:       "payment-" + userID + "@qa.team",
		Password:    "password",
		Gender:      models.Male,
		DateOfBirth: time.Now().AddDate(-20, 0, 0),
		Role:        models.UserRole,
	}
	if err := db.Postgresql.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	event := models.Event{
		ID:          eventID,
		Title:       "Payment Event " + eventID,
		Description: "test event",
		StartDate:   "2026-01-01",
		EndDate:     "2026-01-02",
		Location:    "Lagos",
		TicketPrice: 100,
		Capacity:    100,
		OrganizerID: userID,
	}
	if err := db.Postgresql.Create(&event).Error; err != nil {
		t.Fatal(err)
	}

	ticket := models.Ticket{
		ID:           ticketID,
		EventID:      eventID,
		UserID:       userID,
		PurchaseTime: time.Now(),
		Status:       models.TicketStatusPending,
		Amount:       1,
		Type:         models.TicketTypeRegular,
	}
	if err := db.Postgresql.Create(&ticket).Error; err != nil {
		t.Fatal(err)
	}

	payment := models.Payments{
		ID:            paymentID,
		TicketID:      ticketID,
		PaymentMethod: "card",
		Amount:        100,
		PaymentStatus: models.PendingPayment,
		PaymentTime:   time.Now(),
	}
	if err := db.Postgresql.Create(&payment).Error; err != nil {
		t.Fatal(err)
	}

	return paymentID, ticketID
}

func setupPaymentRouter(t *testing.T) (*gin.Engine, *storage.Database) {
	if os.Getenv("DB_HOST") == "" || os.Getenv("DB_PORT") == "" || os.Getenv("USERNAME") == "" || os.Getenv("DB_NAME") == "" {
		t.Skip("Skipping payment integration tests: DB environment variables are not configured")
	}

	logger := tests.Setup()
	gin.SetMode(gin.TestMode)

	validatorRef := validator.New()
	db := storage.Connection()
	controller := paymentController.Controller{Db: db, Validator: validatorRef, Logger: logger}

	r := gin.Default()
	r.GET("/api/v1/payments/:id", controller.GetPaymentByID)
	r.PATCH("/api/v1/payments/:id/status", controller.UpdatePaymentStatus)
	return r, db
}

func TestGetPaymentByID(t *testing.T) {
	r, db := setupPaymentRouter(t)
	paymentID, _ := seedTicketAndPayment(t, db)

	requestURI := url.URL{Path: "/api/v1/payments/" + paymentID}
	req, err := http.NewRequest(http.MethodGet, requestURI.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	tests.AssertStatusCode(t, rr.Code, http.StatusOK)
	data := tests.ParseResponse(rr)
	statusCode := int(data["status_code"].(float64))
	tests.AssertStatusCode(t, statusCode, http.StatusOK)
}

func TestUpdatePaymentStatusUpdatesTicket(t *testing.T) {
	r, db := setupPaymentRouter(t)
	paymentID, ticketID := seedTicketAndPayment(t, db)

	requestURI := url.URL{Path: "/api/v1/payments/" + paymentID + "/status"}
	body := models.UpdatePaymentStatusRequest{Status: models.SuccessfulPayment}

	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)

	req, err := http.NewRequest(http.MethodPatch, requestURI.String(), &b)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	tests.AssertStatusCode(t, rr.Code, http.StatusOK)

	var ticket models.Ticket
	ticket.ID = ticketID
	if err := ticket.GetTicketByID(db.Postgresql); err != nil {
		t.Fatal(err)
	}

	if ticket.Status != models.TicketStatusApproved {
		t.Fatalf("expected ticket status %d, got %d", models.TicketStatusApproved, ticket.Status)
	}
}
