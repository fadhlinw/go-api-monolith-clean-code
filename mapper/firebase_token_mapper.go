package mapper

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
)

func ToFirebaseTokenModel(dto dto.FirebaseTokenRequest) models.FirebaseToken {
	return models.FirebaseToken{
		UserID: dto.UserID,
		Token:  dto.Token,
	}
}

func ToFirebaseTokenResponse(firebaseToken models.FirebaseToken) dto.FirebaseTokenResponse {
	return dto.FirebaseTokenResponse{
		ID:        firebaseToken.ID,
		UserID:    firebaseToken.UserID,
		Token:     firebaseToken.Token,
		CreatedAt: firebaseToken.CreatedAt,
		UpdatedAt: firebaseToken.UpdatedAt,
	}
}

func ToFirebaseTokensResponse(firebaseTokens []models.FirebaseToken) []dto.FirebaseTokenResponse {
	var firebaseTokensResponse = []dto.FirebaseTokenResponse{}
	for _, firebaseToken := range firebaseTokens {
		firebaseTokensResponse = append(firebaseTokensResponse, ToFirebaseTokenResponse(firebaseToken))
	}
	return firebaseTokensResponse
}
