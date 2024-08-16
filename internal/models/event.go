package models

import "gorm.io/gorm"

type Event struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	Title       string  `json:"title" gorm:"size:255;not null" `
	Description string  `json:"description" gorm:"not null"`
	StartDate   string  `json:"start_date" gorm:"not null"`
	EndDate     string  `json:"end_date" gorm:"not null"`
	Location    string  `json:"location" gorm:"size:255"`
	TicketPrice float64 `json:"ticket_price"`
	Capacity    int     `json:"capacity" gorm:"not null"`
	OrganizerID string  `json:"organizer_id" gorm:"not null"`
	Organizer   User    `json:"organizer" gorm:"foreignKey:OrganizerID"`
	gorm.Model
}

type CreateEventRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	TicketPrice float64 `json:"ticket_price" binding:"required"`
	Capacity    int     `json:"capacity" binding:"required"`
	OrganizerID string  `json:"organizer_id" binding:"required"`
}

type CreateEventResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Location    string  `json:"location"`
	TicketPrice float64 `json:"ticket_price"`
	Capacity    int     `json:"capacity"`
	OrganizerID string  `json:"organizer_id"`
}

type GetEventResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Location    string  `json:"location"`
	TicketPrice float64 `json:"ticket_price"`
	Capacity    int     `json:"capacity"`
	OrganizerID string  `json:"organizer_id"`
}
type GetEventRequest struct {
	ID string `json:"id"`
}
type UpdateEventRequest struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     string  `json:"end_date,omitempty"`
	Location    string  `json:"location,omitempty"`
	TicketPrice float64 `json:"ticket_price,omitempty"`
	Capacity    int     `json:"capacity,omitempty"`
	OrganizerID string  `json:"organizer_id,omitempty"`
}

func (event *Event) CreateEvent(db *gorm.DB) error {
	if err := db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (event *Event) GetEventByID(db *gorm.DB) error {
	if err := db.Where("id = ?", event.ID).First(&event).Error; err != nil {
		return err
	}
	return nil
}

func (event *Event) UpdateEvent(db *gorm.DB) error {
	if err := db.Save(&event).Error; err != nil {
		return err
	}
	return nil
}

func (event *Event) DeleteEvent(db *gorm.DB) error {
	if err := db.Delete(&event).Error; err != nil {
		return err
	}
	return nil
}
