package models

import "time"

type TokenStore struct {
	ID        int
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t TokenStore) TableName() string {
	return "token_store"
}
