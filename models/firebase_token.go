package models

import "time"

type FirebaseToken struct {
	ID        int
	UserID    int
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f FirebaseToken) TableName() string {
	return "firebase_tokens"
}
