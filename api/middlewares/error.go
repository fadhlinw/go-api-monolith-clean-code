package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gitlab.com/tsmdev/software-development/backend/go-project/error"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gorm.io/gorm"
)

type ErrorMiddleware struct {
	logger lib.Logger
}

func NewErrorMiddleware(
	logger lib.Logger,
) ErrorMiddleware {
	return ErrorMiddleware{
		logger: logger,
	}
}

func (m ErrorMiddleware) Setup() {}

func (m ErrorMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case error.Http:
				abortErrorResponse(c, e.StatusCode, e.Description)
			case *mysql.MySQLError:
				if e.Number == 1062 {
					abortErrorResponse(c, http.StatusConflict, e.Message)
				} else {
					abortErrorResponse(c, http.StatusInternalServerError, e.Message)
				}
			default:
				if e == gorm.ErrRecordNotFound {
					abortErrorResponse(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				} else {
					abortErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
			}
		}
	}
}
