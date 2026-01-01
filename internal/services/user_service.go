// services/user_service.go
package services

import (
	"Visa/internal/repos"
	"Visa/models"
	"Visa/pkg"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repos.UserRepo
}

func NewUserService(userRepo *repos.UserRepo) *UserService {
	return &UserService{Repo: userRepo}
}

// Login authenticates a user and returns a JWT token
func (us *UserService) Login(email, password string) (*models.User, string, error) {
	// Validate input
	if email == "" {
		return nil, "", errors.New("email is required")
	}
	if password == "" {
		return nil, "", errors.New("password is required")
	}

	// Get user by email
	user, err := us.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Compare hashed password with plain text password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token using pkg.GenerateToken (correct one!)
	token, err := pkg.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Don't send password hash to frontend
	user.Password = ""

	return user, token, nil
}

// GetAllUsers retrieves all users
func (us *UserService) GetAllUsers() ([]models.User, error) {
	users, err := us.Repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// CreateUser creates a new user with password hashing and validation
func (us *UserService) CreateUser(user *models.User) error {
	if user == nil {
		return errors.New("user data is required")
	}

	// Validate user data
	if err := us.validateUser(user); err != nil {
		return err
	}

	// Check if email already exists
	exists, err := us.Repo.EmailExists(user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return errors.New("email already registered")
	}

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	if err := us.Repo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password before returning
	user.Password = ""
	return nil
}

// GetUserById retrieves a user by their ID
func (us *UserService) GetUserById(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := us.Repo.GetUserById(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Don't send password
	user.Password = ""
	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Validate email format
	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	user, err := us.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// UpdateUser updates an existing user
func (us *UserService) UpdateUser(user *models.User) error {
	if user == nil {
		return errors.New("user data is required")
	}
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	// Verify user exists
	existingUser, err := us.Repo.GetUserById(user.ID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Validate updated fields if provided
	if user.Name != "" {
		if len(user.Name) < 2 {
			return errors.New("name must be at least 2 characters")
		}
		if len(user.Name) > 100 {
			return errors.New("name must not exceed 100 characters")
		}
		existingUser.Name = user.Name
	}

	if user.Email != "" && user.Email != existingUser.Email {
		if !isValidEmail(user.Email) {
			return errors.New("invalid email format")
		}

		// Check if new email is already taken
		exists, err := us.Repo.EmailExists(user.Email)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return errors.New("email already registered")
		}
		existingUser.Email = user.Email
	}

	// Hash new password if provided
	if user.Password != "" {
		if len(user.Password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
		if len(user.Password) > 72 {
			return errors.New("password must not exceed 72 characters")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		existingUser.Password = string(hashedPassword)
	}

	if err := us.Repo.UpdateUser(existingUser); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user
func (us *UserService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("invalid user ID")
	}

	// Verify user exists
	_, err := us.Repo.GetUserById(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := us.Repo.DeleteUser(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ChangePassword changes a user's password
func (us *UserService) ChangePassword(userId uint, oldPassword, newPassword string) error {
	if userId == 0 {
		return errors.New("invalid user ID")
	}
	if oldPassword == "" {
		return errors.New("old password is required")
	}
	if newPassword == "" {
		return errors.New("new password is required")
	}
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}
	if len(newPassword) > 72 {
		return errors.New("new password must not exceed 72 characters")
	}

	// Get user
	user, err := us.Repo.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := us.Repo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// GetUsersByRole retrieves all users with a specific role
func (us *UserService) GetUsersByRole(role string) ([]models.User, error) {
	if role == "" {
		return nil, errors.New("role is required")
	}

	// Validate role
	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
	}
	if !validRoles[role] {
		return nil, errors.New("invalid role. Must be 'user' or 'admin'")
	}

	users, err := us.Repo.GetUsersByRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by role: %w", err)
	}

	// Remove passwords
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

// validateUser validates user data
func (us *UserService) validateUser(user *models.User) error {
	// Validate name
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}
	if len(user.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if len(user.Name) > 100 {
		return errors.New("name must not exceed 100 characters")
	}

	// Validate email
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is required")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if len(user.Password) > 72 {
		return errors.New("password must not exceed 72 characters")
	}

	// Validate role
	if user.Role != "" {
		validRoles := map[string]bool{
			"user":  true,
			"admin": true,
		}
		if !validRoles[user.Role] {
			return errors.New("invalid role. Must be 'user' or 'admin'")
		}
	}

	return nil
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
