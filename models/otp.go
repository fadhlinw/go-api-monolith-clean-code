package models

import "time"

type Otp struct {
	ID        int
	Code      string
	IsUsed    bool
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Otp) TableName() string {
	return "otp"
}
