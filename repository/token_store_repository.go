package repository

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

// TokenStoreRepository handles token store operations
type TokenStoreRepository struct {
	lib.Database
	logger lib.Logger
}

// NewTokenStoreRepository creates a new token store repository
func NewTokenStoreRepository(db lib.Database, logger lib.Logger) TokenStoreRepository {
	return TokenStoreRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r TokenStoreRepository) WithTrx(trxHandle *gorm.DB) TokenStoreRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context.")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

// CreateToken saves a new token in the database
func (r TokenStoreRepository) CreateToken(tokenStore *models.TokenStore) error {
	return r.Database.Create(tokenStore).Error
}

// DeleteToken deletes a token by its value
func (r TokenStoreRepository) DeleteToken(token string) error {
	return r.Database.Where("token = ?", token).Delete(&models.TokenStore{}).Error
}

// ValidateToken checks if a token exists in the database
func (r TokenStoreRepository) ValidateToken(token string) error {
	return r.Database.Where("token = ?", token).First(&models.TokenStore{}).Error
}
