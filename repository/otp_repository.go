package repository

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

// OTPRepository handles database operations for OTP
type OTPRepository struct {
	lib.Database
	logger lib.Logger
}

// NewOTPRepository creates a new instance of OTPRepository
func NewOTPRepository(db lib.Database, logger lib.Logger) OTPRepository {
	return OTPRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository to work with a transaction
func (r OTPRepository) WithTrx(trxHandle *gorm.DB) OTPRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context.")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

// CreateOTP saves a new OTP to the database
func (r OTPRepository) CreateOTP(otp *models.Otp) error {
	return r.Database.Create(otp).Error
}

// UpdateOTPById updates the OTP's is_used field by ID
func (r OTPRepository) UpdateOTPById(id int, isUsed bool) error {
	return r.Database.Model(&models.Otp{}).Where("id = ?", id).Update("is_used", isUsed).Error
}

// GetOTPByCode retrieves an OTP from the database by user ID and code
func (r OTPRepository) GetOTPByCode(userId int, code string) (*models.Otp, error) {
	var otp models.Otp
	err := r.Database.Where("user_id = ? AND code = ? AND is_used = false", userId, code).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// GetOTPByUserIdAndIsUsed retrieves an OTP from the database by user ID and is_used status
func (r OTPRepository) GetOTPByUserIdAndIsUsed(userId int, isUsed bool) (*models.Otp, error) {
	var otp models.Otp
	err := r.Database.Where("user_id = ? AND is_used = ?", userId, isUsed).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}
