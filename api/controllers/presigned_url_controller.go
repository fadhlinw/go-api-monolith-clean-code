package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
)

// DeviceController data type
type AWSPresignedController struct {
	logger       lib.Logger
	awsPresigned lib.S3Uploader
}

// NewDeviceController creates new device controller
func NewAWSPresignedController(logger lib.Logger, awsPresigned lib.S3Uploader) AWSPresignedController {
	return AWSPresignedController{
		logger:       logger,
		awsPresigned: awsPresigned,
	}
}

// GetOneDevice gets one device model
func (d AWSPresignedController) GeneratePreSignedURL(c *gin.Context) {

	presignedReq := dto.AWSPresignedRequest{}

	err := c.ShouldBindJSON(&presignedReq)
	if err != nil {
		d.logger.Error(err)
		c.Error(err)
		return
	}

	presignedURL, err := d.awsPresigned.GeneratePreSignedURL(presignedReq.FileName)
	if err != nil {
		d.logger.Error(err)
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		presignedURL,
		"Pre-signed URL generated successfully",
	)
}
