package domains

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

type AuthService interface {
	WithTrx(trxHandle *gorm.DB) AuthService
	Authorize(tokenString string, claimType string) (*dto.AuthIdentityDto, error)
	CreateToken(user *models.User, claimType string) (string, error)
	CreateRefreshToken(user models.User) (string, error)
	ValidateAuth(request *dto.AuthRequestDto, refreshTokenString string) (*dto.AuthResponseDto, error)
	SaveToken(token string) error
	ChangePassword(id int, request *dto.AuthChangePasswordDto) error
	ForgotPassword(request dto.AuthForgotPasswordDto) error
	ValidateOTP(request dto.ValidateOTPRequestDto) (*dto.ValidateOTPResponseDto, error)
	ResetPassword(request dto.AuthResetPasswordDto) error
	Logout(tokenString string, refreshToken string) error
}
