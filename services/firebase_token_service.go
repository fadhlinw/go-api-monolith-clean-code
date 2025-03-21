package services

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/mapper"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gitlab.com/tsmdev/software-development/backend/go-project/repository"
)

type FirebaseTokenService struct {
	repository repository.FirebaseTokenRepository
}

func NewFirebaseTokenService(repository repository.FirebaseTokenRepository) FirebaseTokenService {
	return FirebaseTokenService{
		repository: repository,
	}
}

func (s FirebaseTokenService) SaveFirebaseToken(userID int, token string) error {
	firebaseToken := models.FirebaseToken{
		UserID: userID,
		Token:  token,
	}
	return s.repository.Save(firebaseToken)
}

func (s FirebaseTokenService) GetFirebaseTokenByUserID(userID int) (*dto.FirebaseTokenResponse, error) {
	firebaseToken, err := s.repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	resp := mapper.ToFirebaseTokenResponse(*firebaseToken)
	return &resp, nil
}

func (s FirebaseTokenService) DeleteToken(userID int) error {
	return s.repository.DeleteByID(userID)
}
