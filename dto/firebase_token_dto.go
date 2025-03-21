package dto

import "time"

type FirebaseTokenRequest struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token" binding:"required"`
}

type FirebaseTokenResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
