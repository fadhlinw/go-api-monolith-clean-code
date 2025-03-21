package routes

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/api/controllers"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/middlewares"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
)

// DeviceModelRoutes struct
type AWSPresignedRoutes struct {
	logger                 lib.Logger
	handler                lib.RequestHandler
	awsPresignedController controllers.AWSPresignedController
	authMiddleware         middlewares.JWTAuthMiddleware
	errorMiddleware        middlewares.ErrorMiddleware
}

// Setup device model routes
func (r AWSPresignedRoutes) Setup() {
	r.logger.Info("Setting up routes")
	api := r.handler.Gin.Group("/api").Use(r.authMiddleware.Handler(), r.errorMiddleware.Handler())
	{
		api.POST("/presigned-url", r.awsPresignedController.GeneratePreSignedURL)
	}
}

// NewDeviceModelRoutes creates new device model controller
func NewAWSPresignedRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	awsPresignedController controllers.AWSPresignedController,
	authMiddleware middlewares.JWTAuthMiddleware,
	errorMiddleware middlewares.ErrorMiddleware,
) AWSPresignedRoutes {
	return AWSPresignedRoutes{
		handler:                handler,
		logger:                 logger,
		awsPresignedController: awsPresignedController,
		authMiddleware:         authMiddleware,
		errorMiddleware:        errorMiddleware,
	}
}
