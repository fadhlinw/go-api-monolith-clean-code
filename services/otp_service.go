package services

import (
	"net/http"

	"gitlab.com/tsmdev/software-development/backend/go-project/constants"
	"gitlab.com/tsmdev/software-development/backend/go-project/domains"
	httperror "gitlab.com/tsmdev/software-development/backend/go-project/error"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gitlab.com/tsmdev/software-development/backend/go-project/repository"
)

// OTPService handles OTP-related logic
type OTPService struct {
	env        lib.Env
	logger     lib.Logger
	repository repository.OTPRepository
}

// NewOTPService creates a new OTPService instance
func NewOTPService(env lib.Env, logger lib.Logger, repository repository.OTPRepository) domains.OTPService {
	return OTPService{
		env:        env,
		logger:     logger,
		repository: repository,
	}
}

// Create saves the OTP data to the database
func (s OTPService) Create(userId int, code string) error {
	otp := models.Otp{
		UserId: userId,
		Code:   code,
		IsUsed: false,
	}
	err := s.repository.CreateOTP(&otp)
	if err != nil {
		s.logger.Error("Error saving OTP data to DB")
		s.logger.Debug("Detail: ", err.Error())
		return httperror.NewHttpError(constants.ERROR_CREATING_OTP, "", http.StatusInternalServerError)
	}
	return nil
}

// UpdateById updates the OTP data in the database by ID
func (s OTPService) UpdateById(id int, isUsed bool) error {
	s.logger.Info("Updating OTP data in DB")
	err := s.repository.UpdateOTPById(id, isUsed)
	if err != nil {
		s.logger.Error("Error updating OTP data to DB")
		s.logger.Debug("Detail: ", err.Error())
		return httperror.NewHttpError(constants.ERROR_UPDATING_OTP, "", http.StatusInternalServerError)
	}
	return nil
}

// GetByCode retrieves the OTP data from the database by code
func (s OTPService) GetByCode(userId int, code string) (*models.Otp, error) {
	s.logger.Info("Getting OTP data from DB")
	otp, err := s.repository.GetOTPByCode(userId, code)
	if err != nil {
		s.logger.Error("Error getting OTP data from DB")
		s.logger.Debug("Detail: ", err.Error())
		return nil, httperror.NewHttpError(constants.ERROR_GETTING_OTP_BY_CODE, "", http.StatusInternalServerError)
	}
	return otp, nil
}

// GetByUserIdAndIsUsed retrieves the OTP data from the database by user ID and isUsed flag
func (s OTPService) GetByUserIdAndIsUsed(userId int, isUsed bool) (*models.Otp, error) {
	s.logger.Info("Getting OTP data from DB")
	otp, err := s.repository.GetOTPByUserIdAndIsUsed(userId, isUsed)
	if err != nil {
		s.logger.Error("Error getting OTP data from DB")
		s.logger.Debug("Detail: ", err.Error())
		return nil, err
	}
	return otp, nil
}
