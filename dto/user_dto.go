package dto

import "time"

type UserRequestDto struct {
}

type CreateUserRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Age          uint8  `json:"age" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
	Birthday     string `json:"birthday" binding:"required"`
	MemberNumber string `json:"member_number" binding:"required"`
}

type UserResponseDto struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Age          uint8     `json:"age"`
	Gender       string    `json:"gender"`
	Birthday     string    `json:"birthday"`
	MemberNumber string    `json:"member_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
