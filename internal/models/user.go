package models

import "gorm.io/gorm"
const (
	AdminRole = "admin"
	UserRole  = "user"
	OrganizerRole = "organizer"
)

type User struct {
	ID       string    `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role    string `json:"role"`
	gorm.Model
}