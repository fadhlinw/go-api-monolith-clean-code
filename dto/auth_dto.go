package dto

type AuthRequestDto struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Client    string `json:"client"`    // temporary optional field
	Signature string `json:"signature"` // temporary optional field
}

type AuthIdentityDto struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
}

type AuthIssueRefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponseDto struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthChangePasswordDto struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type AuthForgotPasswordDto struct {
	Email string `json:"email" binding:"required"`
}

type AuthResetPasswordDto struct {
	UserId     int    `json:"user_id"`
	Password   string `json:"password" binding:"required"`
	ResetToken string `json:"reset_token"`
}

type LogoutRequestDto struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
