package services

import (
	"math"

	"gitlab.com/tsmdev/software-development/backend/go-project/utils"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module exports services present
var Module = fx.Options(
	fx.Provide(NewUserService),
	fx.Provide(NewJWTAuthService),
	fx.Provide(NewOTPService),
	fx.Provide(NewTokenStoreService),
)

func paginate(value interface{}, pagination *utils.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
