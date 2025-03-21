package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/tsmdev/software-development/backend/go-project/constants"
	"gitlab.com/tsmdev/software-development/backend/go-project/domains"
	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/utils"
	"gitlab.com/tsmdev/software-development/backend/go-project/validation"

	"gorm.io/gorm"
)

// UserController data type
type UserController struct {
	service domains.UserService
	logger  lib.Logger
}

// NewUserController creates new user controller
func NewUserController(userService domains.UserService, logger lib.Logger) UserController {
	return UserController{
		service: userService,
		logger:  logger,
	}
}

// GetOneUser gets one user
func (u UserController) GetOneUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	user, err := u.service.GetOneUserById(id)
	if err != nil {
		u.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.UserResponseDto{}),
			nil,
			"Bad request",
		)
		return
	}
	globalResponse(
		c,
		http.StatusOK,
		nil,
		user,
		"Success",
	)

}

// GetUser gets the user
func (u UserController) GetUser(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	sort := c.DefaultQuery("sort", "created_at,desc")
	searchQuery := c.DefaultQuery("q", "")
	pagination := utils.Pagination{
		Page:  page + 1,
		Limit: limit,
		Sort:  utils.ReplaceComaWithSpace(sort),
	}

	users, err := u.service.GetAllUser(searchQuery, pagination)
	if err != nil {
		u.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.UserResponseDto{}),
			nil,
			"Bad request",
		)
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		users,
		"Success",
	)

}

// SaveUser saves the user
func (u UserController) SaveUser(c *gin.Context) {
	user := dto.CreateUserRequest{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.CreateUserRequest{}),
			nil,
			"Bad request",
		)
		return
	}

	if err := u.service.WithTrx(trxHandle).CreateUser(user); err != nil {
		u.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.CreateUserRequest{}),
			nil,
			"Bad request",
		)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		"User created successfully",
	)
}

// UpdateUser updates user
func (u UserController) UpdateUser(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.logger.Error(err)
		c.Error(err)
		return
	}

	user := dto.CreateUserRequest{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)
	if err := c.ShouldBindJSON(&user); err != nil {
		u.logger.Error(err)
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.CreateUserRequest{}),
			nil,
			"Bad request",
		)
		return
	}

	if err := u.service.WithTrx(trxHandle).UpdateUser(id, user); err != nil {
		u.logger.Error(err)
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		"User updated successfully",
	)
}

// DeleteUser deletes user
func (u UserController) DeleteUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		globalResponse(
			c,
			http.StatusBadRequest,
			validation.ParseFieldErrors(err, dto.CreateUserRequest{}),
			nil,
			"Bad request",
		)
		return
	}

	if err := u.service.DeleteUser(id); err != nil {
		u.logger.Error(err)
		c.Error(err)
		return
	}

	globalResponse(
		c,
		http.StatusOK,
		nil,
		nil,
		"User deleted successfully",
	)
}
