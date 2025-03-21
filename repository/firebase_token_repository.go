package repository

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

// FirebaseTokenRepository handles database operations for Firebase Tokens
type FirebaseTokenRepository struct {
	lib.Database
	logger lib.Logger
}

// NewFirebaseTokenRepository creates a new FirebaseTokenRepository
func NewFirebaseTokenRepository(db lib.Database, logger lib.Logger) FirebaseTokenRepository {
	return FirebaseTokenRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r FirebaseTokenRepository) WithTrx(trxHandle *gorm.DB) FirebaseTokenRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context.")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

// GetByUserID retrieves a Firebase Token by user ID
func (r FirebaseTokenRepository) GetByUserID(userID int) (*models.FirebaseToken, error) {
	token := models.FirebaseToken{}
	err := r.Database.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetAll retrieves all Firebase Tokens
func (r FirebaseTokenRepository) GetAll() ([]models.FirebaseToken, error) {
	firebaseTokens := []models.FirebaseToken{}
	err := r.Database.Find(&firebaseTokens).Error
	return firebaseTokens, err
}

// Save saves or updates a Firebase Token
func (r FirebaseTokenRepository) Save(firebaseToken models.FirebaseToken) error {
	return r.Database.Save(firebaseToken).Error
}

// DeleteByID deletes a Firebase Token by ID
func (r FirebaseTokenRepository) DeleteByID(id int) error {
	return r.Database.Where("id = ?", id).Delete(&models.FirebaseToken{}).Error
}
