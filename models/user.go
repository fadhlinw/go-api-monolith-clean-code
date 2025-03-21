package models

import (
	"time"
)

// User model
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Gender       string    `json:"gender"`
	Age          uint8     `json:"age"`
	Birthday     string    `json:"birthday"`
	MemberNumber string    `json:"member_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return "users"
}
