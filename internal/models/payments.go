package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	SuccessfulPayment = iota
	FailedPayment
)
// Payment model struct

type Payments struct {
    ID            string           `json:"id" gorm:"primaryKey"`                  // id SERIAL PRIMARY KEY
    TicketID      string           `json:"ticket_id"`                             // ticket_id  REFERENCES tickets(id)
    Ticket        Ticket         `json:"ticket" gorm:"foreignKey:TicketID"`     // Foreign key relationship with Ticket
    PaymentMethod string         `json:"payment_method" gorm:"size:50;not null"`// payment_method VARCHAR(50) NOT NULL
    Amount        float64        `json:"amount" gorm:"type:decimal(10,2);not null"` // amount DECIMAL(10, 2) NOT NULL
    PaymentStatus int        `json:"payment_status" gorm:"size:50;not null"`// payment_status VARCHAR(50) NOT NULL CHECK (payment_status IN ('successful', 'failed'))
    PaymentTime   time.Time      `json:"payment_time" gorm:"autoCreateTime"`    // payment_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    gorm.Model

}