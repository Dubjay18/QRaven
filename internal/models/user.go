package models

import "gorm.io/gorm"

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role    string `json:"role"`
	gorm.Model
}