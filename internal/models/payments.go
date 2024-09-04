package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	SuccessfulPayment = iota
	PendingPayment
	FailedPayment
)

// Payment model struct

type Payments struct {
	ID            string    `json:"id" gorm:"primaryKey"`                      // id SERIAL PRIMARY KEY
	TicketID      string    `json:"ticket_id"`                                 // ticket_id  REFERENCES tickets(id)
	Ticket        Ticket    `json:"ticket" gorm:"foreignKey:TicketID"`         // Foreign key relationship with Ticket
	PaymentMethod string    `json:"payment_method" gorm:"size:50;not null"`    // payment_method VARCHAR(50) NOT NULL
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"` // amount DECIMAL(10, 2) NOT NULL
	PaymentStatus int       `json:"payment_status" gorm:"size:50;not null"`    // payment_status VARCHAR(50) NOT NULL CHECK (payment_status IN ('successful', 'failed'))
	PaymentTime   time.Time `json:"payment_time" gorm:"autoCreateTime"`        // payment_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	gorm.Model
}
type PaystackInitializeRequest struct {
	Email     string `json:"email"`
	Amount    int    `json:"amount"` // Amount in kobo (1 Naira = 100 kobo)
	Reference string `json:"reference"`
	Currency  string `json:"currency"`
}

type PaystackInitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type InitializePaymentRequest struct {
	TicketID      string  `json:"ticket_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Email         string  `json:"email" binding:"required"`
}

type InitializePaymentResponse struct {
	ID            string  `json:"id"`
	TicketID      string  `json:"ticket_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	PaymentStatus int     `json:"payment_status"`
	PaymentTime   string  `json:"payment_time"`
}

func (p *Payments) Create(db *gorm.DB) error {
	return db.Create(p).Error
}

func (p *Payments) GetPaymentsByTicketID(db *gorm.DB, ticketID string) ([]Payments, error) {
	var payments []Payments
	if err := db.Where("ticket_id = ?", ticketID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (p *Payments) GetPaymentByID(db *gorm.DB, id string) (*Payments, error) {
	var payment Payments
	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *Payments) UpdatePaymentStatus(db *gorm.DB, id string, status int) error {
	return db.Model(&Payments{}).Where("id = ?", id).Update("payment_status", status).Error
}
