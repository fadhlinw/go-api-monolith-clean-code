package routes

import (
	"gitlab.com/tsmdev/software-development/backend/go-project/api/controllers"
	"gitlab.com/tsmdev/software-development/backend/go-project/api/middlewares"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
)

// UserRoutes struct
type UserRoutes struct {
	logger          lib.Logger
	handler         lib.RequestHandler
	userController  controllers.UserController
	authMiddleware  middlewares.JWTAuthMiddleware
	errorMiddleware middlewares.ErrorMiddleware
}

// Setup user routes
func (s UserRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.Handler(), s.errorMiddleware.Handler())
	{
		api.GET("/users", s.userController.GetUser)
		api.GET("/users/:id", s.userController.GetOneUser)
		api.POST("/users", s.userController.SaveUser)
		api.POST("/users/:id", s.userController.UpdateUser)
		api.DELETE("/users/:id", s.userController.DeleteUser)
	}
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	userController controllers.UserController,
	authMiddleware middlewares.JWTAuthMiddleware,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
		authMiddleware: authMiddleware,
	}
}
