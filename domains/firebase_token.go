package domains

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

type FirebaseTokenService interface {
	WithTrx(trxHandle *gorm.DB) FirebaseTokenService
	GetFirebaseTokensByUserID(userIDs []uint) ([]models.FirebaseToken, error)
	GetFirebaseTokenByUserID(userID uint) (*models.FirebaseToken, error)
	GetAllFirebaseToken() ([]dto.FirebaseTokenResponse, error)
	SaveTokenByUserId(request dto.FirebaseTokenRequest) error
	DeleteFirebaseTokenById(id int) error
}
