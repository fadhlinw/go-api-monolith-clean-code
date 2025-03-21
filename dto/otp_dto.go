package dto

type ValidateOTPRequestDto struct {
	Code  string `json:"code" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type ValidateOTPResponseDto struct {
	ResetToken string `json:"reset_token"`
}
