package models

import "gorm.io/gorm"

type Notification struct {
    ID        string           `json:"id" gorm:"primaryKey"`               // id SERIAL PRIMARY KEY
    UserID    string           `json:"user_id"`                            // user_id  REFERENCES users(id)
    User      User           `json:"user" gorm:"foreignKey:UserID"`      // Foreign key relationship with User
    Message   string         `json:"message" gorm:"type:text;not null"`  // message TEXT NOT NULL
    Read      bool           `json:"read" gorm:"default:false"`          // read BOOLEAN DEFAULT FALSE
	gorm.Model
}