package models

import "gorm.io/gorm"

type Event struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"not null"`
	StartDate   string `json:"start_date" gorm:"not null"`
	EndDate     string `json:"end_date" gorm:"not null"`
	Location   string `json:"location" gorm:"size:255"`
	TicketPrice   float64 `json:"ticket_price"`
	Capacity	int    `json:"capacity" gorm:"not null"`
	OrganizerID string    `json:"organizer_id" gorm:"not null"`
	Organizer   User `json:"organizer" gorm:"foreignKey:OrganizerID"`
	gorm.Model
}
