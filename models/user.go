package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username string    `gorm:"unique;not null" json:"username"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	IsAdmin  bool      `gorm:"default:false;not null" json:"isAdmin"`
	Accounts []Account `json:"-"`

	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
}

func NewUser(
	username string,
	email string,
	password string,
) *User {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password,
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
