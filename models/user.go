package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToUser(user database.User) *User {
	return &User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
