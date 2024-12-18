package models

import "gorm.io/gorm"

type Event struct {
	ID          string  `json:"id" gorm:"primaryKey"`
	Title       string  `json:"title" gorm:"size:255;not null" `
	Image       string  `json:"image"`
	Description string  `json:"description" gorm:"not null"`
	StartDate   string  `json:"start_date" gorm:"not null"`
	EndDate     string  `json:"end_date" gorm:"not null"`
	Location    string  `json:"location" gorm:"size:255"`
	TicketPrice float64 `json:"ticket_price"`
	Capacity    int     `json:"capacity" gorm:"not null"`
	OrganizerID string  `json:"organizer_id" gorm:"not null"`
	Organizer   User    `json:"organizer" gorm:"foreignKey:OrganizerID"`
	TicketCount int64   `json:"ticket_count" gorm:"default:0"`
	gorm.Model
}

type CreateEventRequest struct {
	Title       string  `form:"title" binding:"required"`
	Image       string  `form:"image"`
	Description string  `form:"description" binding:"required"`
	StartDate   string  `form:"start_date" binding:"required"`
	EndDate     string  `form:"end_date" binding:"required"`
	Location    string  `form:"location" binding:"required"`
	TicketPrice float64 `form:"ticket_price" binding:"required"`
	Capacity    int     `form:"capacity" binding:"required"`
	OrganizerID string  `form:"organizer_id" binding:"required"`
}

type CreateEventResponse struct {
	ID          string  `json:"id"`
	Image       string  `json:"image"`
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
	Image       string  `json:"image"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	Location    string  `json:"location"`
	TicketPrice float64 `json:"ticket_price"`
	Capacity    int     `json:"capacity"`
	OrganizerID string  `json:"organizer_id"`
	TicketCount int64   `json:"ticket_count"`
}
type GetEventRequest struct {
	ID string `json:"id"`
}
type UpdateEventRequest struct {
	Title       string  `json:"title,omitempty"`
	Image       string  `json:"image,omitempty"`
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

func (event *Event) IncreaseTicketCount(db *gorm.DB, amount int) error {
	if err := db.Model(&event).Update("ticket_count", gorm.Expr("ticket_count + ?", amount)).Error; err != nil {
		return err
	}
	return nil
}
