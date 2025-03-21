package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewUserController),
	fx.Provide(NewJWTAuthController),
	fx.Provide(NewAWSPresignedController),
)

func globalResponse(c *gin.Context, statusCode int, errors interface{}, data interface{}, message interface{}) {
	c.JSON(statusCode, gin.H{
		"error":   errors,
		"data":    data,
		"message": message,
	})
}
