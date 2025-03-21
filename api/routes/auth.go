package routes

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/api/controllers"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/middlewares"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
)

// AuthRoutes struct
type AuthRoutes struct {
	logger          lib.Logger
	handler         lib.RequestHandler
	authController  controllers.JWTAuthController
	errorMiddleware middlewares.ErrorMiddleware
	authMiddleware  middlewares.JWTAuthMiddleware
	resetMiddleware middlewares.JWTResetPasswordMiddleware
}

// Setup user routes
func (s AuthRoutes) Setup() {
	s.logger.Info("Setting up routes")

	reset := s.handler.Gin.Group("/api/auth/reset_password").Use(s.resetMiddleware.Handler(), s.errorMiddleware.Handler())
	{
		reset.POST("", s.authController.ResetPassword)
	}

	logout := s.handler.Gin.Group("api/auth/logout").Use(s.authMiddleware.Handler(), s.errorMiddleware.Handler())
	{
		logout.POST("", s.authController.Logout)
	}

	auth := s.handler.Gin.Group("/api/auth").Use(s.errorMiddleware.Handler())
	{
		auth.POST("/login", s.authController.SignIn)
		auth.POST("/v2/login", s.authController.SignInV2)
		auth.POST("/register", s.authController.Register)
		auth.POST("/issue_access_token", s.authController.IssueAccessToken)
		auth.POST("/forgot_password", s.authController.ForgotPassword)
		auth.POST("/validate_otp", s.authController.ValidateOTP)
	}

	privateAuth := s.handler.Gin.Group("/api/private").Use(s.authMiddleware.Handler(), s.errorMiddleware.Handler())
	{
		privateAuth.POST("/auth/change_password", s.authController.ChangePassword)
	}
}

// NewAuthRoutes creates new user controller
func NewAuthRoutes(
	handler lib.RequestHandler,
	authController controllers.JWTAuthController,
	logger lib.Logger,
	errorMiddleware middlewares.ErrorMiddleware,
	authMiddleware middlewares.JWTAuthMiddleware,
	resetMiddleware middlewares.JWTResetPasswordMiddleware,
) AuthRoutes {
	return AuthRoutes{
		handler:         handler,
		logger:          logger,
		authController:  authController,
		errorMiddleware: errorMiddleware,
		authMiddleware:  authMiddleware,
		resetMiddleware: resetMiddleware,
	}
}
