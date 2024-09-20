package models

import (
	"errors"
	"qraven/pkg/repository/storage/postgresql"
	"time"

	"gorm.io/gorm"
)

type RoleName string
type RoleId int

type DefaultIdentity struct {
	User      RoleId
	Admin     RoleId
	Organizer RoleId
}

var RoleIdentity = DefaultIdentity{
	User:      1,
	Organizer: 2,
	Admin:     3,
}

var (
	AdminRole     RoleName = "admin"
	UserRole      RoleName = "user"
	OrganizerRole RoleName = "organizer"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
	Other  Gender = "other"
)

type User struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Gender      Gender    `json:"gender" gorm:"not null"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"not null"`
	Avatar      string    `json:"avatar"`
	Role        RoleName  `json:"role"`

	gorm.Model
}

type CreateUserRequest struct {
	FirstName   string `form:"first_name" binding:"required"`
	LastName    string `form:"last_name" binding:"required"`
	Email       string `form:"email" binding:"required,email"`
	Password    string `form:"password" binding:"required,min=6"`
	Gender      Gender `form:"gender" binding:"required" gorm:"not null"`
	DateOfBirth string `form:"date_of_birth" binding:"required" gorm:"not null"`
	Role        string `form:"role"`
	Avatar      string `json:"avatar"`
}

type CreateUserResponse struct {
	ID          string `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"requirxed,email"`
	Role        string `json:"role"`
	Gender      Gender `json:"gender" gorm:"not null"`
	Avatar      string `json:"avatar"`
	DateOfBirth string `json:"date_of_birth" gorm:"not null"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (c *CreateUserRequest) ParseDateOfBirth() (time.Time, error) {
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, c.DateOfBirth)
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}
	return parsedDate, nil

}

func (c *User) GetRoleName() RoleName {
	return c.Role
}

func (u User) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
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
