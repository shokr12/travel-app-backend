// repos/user_repo.go
package repos

import (
	"Visa/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetAllUsers retrieves all users from the database
func (ur *UserRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := ur.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserById retrieves a user by their ID
func (ur *UserRepo) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email address
func (ur *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (ur *UserRepo) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

// UpdateUser updates an existing user
func (ur *UserRepo) UpdateUser(user *models.User) error {
	return ur.db.Save(user).Error
}

// DeleteUser deletes a user by their ID
func (ur *UserRepo) DeleteUser(id uint) error {
	return ur.db.Delete(&models.User{}, id).Error
}

// GetUsersByRole retrieves all users with a specific role (e.g., "admin", "user")
func (ur *UserRepo) GetUsersByRole(role string) ([]models.User, error) {
	var users []models.User
	if err := ur.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// EmailExists checks if an email already exists in the database
func (ur *UserRepo) EmailExists(email string) (bool, error) {
	var count int64
	if err := ur.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAdminUsers retrieves all users with admin role
func (ur *UserRepo) GetAdminUsers() ([]models.User, error) {
	return ur.GetUsersByRole("admin")
}

// GetRegularUsers retrieves all users with regular user role
func (ur *UserRepo) GetRegularUsers() ([]models.User, error) {
	return ur.GetUsersByRole("user")
}
