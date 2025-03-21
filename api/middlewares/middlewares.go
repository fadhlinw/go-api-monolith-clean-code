package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewCorsMiddleware),
	fx.Provide(NewJWTAuthMiddleware),
	fx.Provide(NewDatabaseTrx),
	fx.Provide(NewErrorMiddleware),
	fx.Provide(NewMiddlewares),
	fx.Provide(NewJWTResetPasswordMiddleware),
)

// IMiddleware middleware interface
type IMiddleware interface {
	Setup()
}

// Middlewares contains multiple middleware
type Middlewares []IMiddleware

// NewMiddlewares creates new middlewares
// Register the middleware that should be applied directly (globally)
func NewMiddlewares(
	corsMiddleware CorsMiddleware,
	dbTrxMiddleware DatabaseTrx,
	jwtAuthMiddleware JWTAuthMiddleware,
	errorMiddlware ErrorMiddleware,
	jwtResetPasswordMiddleware JWTResetPasswordMiddleware,
) Middlewares {
	return Middlewares{
		corsMiddleware,
		dbTrxMiddleware,
		jwtAuthMiddleware,
		errorMiddlware,
		jwtResetPasswordMiddleware,
	}
}

// Setup sets up middlewares
func (m Middlewares) Setup() {
	for _, middleware := range m {
		middleware.Setup()
	}
}

func abortErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error":   http.StatusText(statusCode),
		"data":    nil,
		"message": errorMessage,
	})
}
