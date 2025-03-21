package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/tsmdev/software-development/backend/go-project/constants"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/domains"
)

// JWTAuthMiddleware middleware for jwt authentication
type JWTAuthMiddleware struct {
	service domains.AuthService
	logger  lib.Logger
}

// NewJWTAuthMiddleware creates new jwt auth middleware
func NewJWTAuthMiddleware(
	logger lib.Logger,
	service domains.AuthService,
) JWTAuthMiddleware {
	return JWTAuthMiddleware{
		service: service,
		logger:  logger,
	}
}

// Setup sets up jwt auth middleware
func (m JWTAuthMiddleware) Setup() {}

// Handler handles middleware functionality
func (m JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := m.service.Authorize(authToken, constants.TypeAuthToken)
			if authorized != nil {
				c.Request.Header.Set("user_id", authorized.UserID)
				c.Next()
				return
			}

			m.logger.Error(err)
			abortErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		abortErrorResponse(c, http.StatusUnauthorized, "invalid token")
	}
}
