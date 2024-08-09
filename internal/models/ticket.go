package models

import (
	"time"

	"gorm.io/gorm"
)

const(
	TicketStatusPending = iota
	TicketStatusApproved
	TicketStatusCancelled
)

const(
	TicketTypeRegular = "regular"
	TicketTypeVip = "vip"
	TicketTypePremium = "premium"
)


type Ticket struct {
    ID           string           `json:"id" gorm:"primaryKey"`                
    EventID      uint           `json:"event_id"`                             // event_id INT REFERENCES events(id)
    Event        Event          `json:"event" gorm:"foreignKey:EventID"`      // Foreign key relationship with Event
    UserID       string           `json:"user_id"`                              // user_id INT REFERENCES users(id)
    User         User           `json:"user" gorm:"foreignKey:UserID"`        // Foreign key relationship with User
    QRCode       string         `json:"qr_code" gorm:"unique"`                // qr_code TEXT UNIQUE
    PurchaseTime time.Time      `json:"purchase_time" gorm:"autoCreateTime"`  // purchase_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    Status       int         `json:"status" gorm:"not null"`        // status in ['pending', 'approved', 'cancelled']
    Amount       float64        `json:"amount" gorm:"type:decimal(10,2);not null"` // amount DECIMAL(10, 2) NOT NULL
	Type 	   string         `json:"type" gorm:"size:50;not null"` // type VARCHAR(50) NOT NULL
	gorm.Model
}