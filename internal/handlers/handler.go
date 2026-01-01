// handlers/user_handler.go
package handlers

import (
	"Visa/internal/services"
	"Visa/models"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=100"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8,max=72"`
}

// CreateUser creates a new user account
func (uh *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := uh.UserService.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// LoginUser authenticates a user and returns a JWT token
func (uh *UserHandler) LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please provide valid email and password",
		})
		return
	}

	user, token, err := uh.UserService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_failed",
			"message": "Invalid email or password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})
}

// GetUserById retrieves a user by ID
func (uh *UserHandler) GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "User ID must be a valid number",
		})
		return
	}

	currentUserID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && uint(id) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only view your own profile",
		})
		return
	}

	user, err := uh.UserService.GetUserById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "User not found",
			})
			return
		}
		log.Printf("Error fetching user %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// UpdateUser updates a user's information
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "User ID must be a valid number",
		})
		return
	}

	currentUserID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && uint(id) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only update your own profile",
		})
		return
	}

	user, err := uh.UserService.GetUserById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "User not found",
			})
			return
		}
		log.Printf("Error fetching user %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve user",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := uh.UserService.UpdateUser(user); err != nil {
		log.Printf("Error updating user %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// DeleteUser deletes a user account
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "User ID must be a valid number",
		})
		return
	}

	currentUserID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && uint(id) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only delete your own account",
		})
		return
	}

	if err := uh.UserService.DeleteUser(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "User not found",
			})
			return
		}
		log.Printf("Error deleting user %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllUsers retrieves all users with pagination (admin only)
func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, err := uh.UserService.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching all users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve users",
		})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(users) {
		start = len(users)
	}
	if end > len(users) {
		end = len(users)
	}

	sanitizedUsers := make([]gin.H, len(users[start:end]))
	for i, user := range users[start:end] {
		sanitizedUsers[i] = gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  sanitizedUsers,
		"page":  page,
		"limit": limit,
		"total": len(users),
	})
}
