package domains

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

type TokenStoreService interface {
	WithTrx(trxHandle *gorm.DB) TokenStoreService
	CreateToken(tokenStore models.TokenStore) error
	DeleteToken(token string) error
	ValidateToken(token string) error
}
