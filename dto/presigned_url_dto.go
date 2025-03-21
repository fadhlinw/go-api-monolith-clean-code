package dto

type AWSPresignedRequest struct {
	FileName string `json:"file_name" binding:"required"`
}
