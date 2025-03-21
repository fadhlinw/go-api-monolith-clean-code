package repository

import (
	"net/http"

	httperror "gitlab.com/tsmdev/software-development/backend/go-project/error"
	"gitlab.com/tsmdev/software-development/backend/go-project/lib"
	"gitlab.com/tsmdev/software-development/backend/go-project/models"
	"gorm.io/gorm"
)

// UserRepository for handling database operations
type UserRepository struct {
	lib.Database
	logger lib.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db lib.Database, logger lib.Logger) UserRepository {
	return UserRepository{
		Database: db,
		logger:   logger,
	}
}

// WithTrx enables repository with transaction
func (r UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		r.logger.Error("Transaction Database not found in gin context.")
		return r
	}
	r.Database.DB = trxHandle
	return r
}

// GetByID retrieves a user by their ID
func (r UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.Database.First(&user, "id = ?", id).Error
	return &user, err
}

// GetByEmail retrieves a user by email
func (r UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.Database.First(&user, "email = ?", email).Error

	if err != nil {
		return nil, httperror.NewHttpError("User not found", "", http.StatusNotFound)
	}
	return &user, err
}

// GetByUsername retrieves a user by username
func (r UserRepository) GetByUsername(email string) (*models.User, error) {
	var user models.User
	err := r.Database.First(&user, "username = ?", email).Error

	if err != nil {
		return nil, httperror.NewHttpError("User not found", "", http.StatusNotFound)
	}
	return &user, err
}

// GetAll retrieves all users with optional search query and pagination
func (r UserRepository) GetAll(searchQuery string) (*gorm.DB, []models.User) {
	var users []models.User
	query := r.Database.Table("users")

	if searchQuery != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	return query, users
}

// Create creates a new user
func (r UserRepository) Create(user *models.User) error {
	return r.Database.Create(user).Error
}

// Update updates user data by ID
func (r UserRepository) Update(id int, user *models.User) error {
	return r.Database.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
}

// Delete removes a user by ID
func (r UserRepository) Delete(id int) error {
	return r.Database.Delete(&models.User{}, id).Error
}
