package models

import (
	"qraven/pkg/repository/storage/postgresql"

	"gorm.io/gorm"
)

type RoleName string
type RoleId int



type DefaultIdentity struct {
	User       RoleId
	Admin 	RoleId
	Organizer RoleId
}

var RoleIdentity = DefaultIdentity{
	User:       1,
	Organizer: 2,
	Admin:      3,	
}

var (
	AdminRole     RoleName = "admin"
	UserRole      RoleName = "user"
	OrganizerRole RoleName = "organizer"
)



type User struct {
	ID       string    `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role    string `json:"role"`
	gorm.Model
}


type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`

}

type CreateUserResponse struct {
	ID       string    `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role    string `json:"role"`
}

func (u *User) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetUserByID(db *gorm.DB, id string) (*User, error) {
	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) UpdateUser(db *gorm.DB, updates map[string]interface{}) error {
	return db.Model(u).Updates(updates).Error
}

func (u *User) CheckUserExistence(db *gorm.DB, id string) bool {
	return postgresql.CheckExistsInTable(db, "users", "id = ?", id)
}