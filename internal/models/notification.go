package models

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID      string `json:"id" gorm:"primaryKey"`              // id SERIAL PRIMARY KEY
	UserID  string `json:"user_id"`                           // user_id  REFERENCES users(id)
	User    User   `json:"user" gorm:"foreignKey:UserID"`     // Foreign key relationship with User
	Message string `json:"message" gorm:"type:text;not null"` // message TEXT NOT NULL
	Read    bool   `json:"read" gorm:"default:false"`         // read BOOLEAN DEFAULT FALSE
	gorm.Model
}

type ExpoPushToken struct {
	ID        string    `json:"id" gorm:"primaryKey"`            // id SERIAL PRIMARY KEY
	UserID    string    `json:"user_id"`                         // user_id  REFERENCES users(id)
	User      User      `json:"user" gorm:"foreignKey:UserID"`   // Foreign key relationship with User
	Token     string    `json:"token" gorm:"type:text;not null"` // token TEXT NOT NULL
	ExpiredAt time.Time `json:"expired_at"`
	gorm.Model
}
