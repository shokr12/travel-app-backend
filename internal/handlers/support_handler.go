// handlers/support_handler.go
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

type SupportHandler struct {
	SupportService *services.SupportService
}

func NewSupportHandler(supportService *services.SupportService) *SupportHandler {
	return &SupportHandler{SupportService: supportService}
}

type CreateTicketRequest struct {
	Subject string `json:"subject" binding:"required,min=5,max=200"`
	Message string `json:"message" binding:"required,min=10,max=2000"`
}

type UpdateTicketRequest struct {
	Subject string `json:"subject" binding:"omitempty,min=5,max=200"`
	Message string `json:"message" binding:"omitempty,min=10,max=2000"`
	Status  string `json:"status" binding:"omitempty,oneof=open in_progress resolved closed"`
}

// CreateTicket creates a new support ticket
func (sh *SupportHandler) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	userID := c.GetUint("userId")

	ticket := models.SupportTicket{
		UserID:  uint(userID),
		Subject: req.Subject,
		Message: req.Message,
		Status:  "open",
	}

	if err := sh.SupportService.CreateSupportTicket(&ticket); err != nil {
		log.Printf("Error creating ticket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to create support ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Support ticket created successfully",
		"data":    ticket,
	})
}

// GetTicketById retrieves a support ticket by ID
func (sh *SupportHandler) GetTicketById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Ticket ID must be a valid number",
		})
		return
	}

	ticket, err := sh.SupportService.GetSupportTicketById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Support ticket not found",
			})
			return
		}
		log.Printf("Error fetching ticket %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve support ticket",
		})
		return
	}

	userID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You don't have permission to view this ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

// GetTicketsByUser retrieves all support tickets for a specific user
func (sh *SupportHandler) GetTicketsByUser(c *gin.Context) {
	// 1 Get userId from URL
	userParam := c.Param("userId")
	userIdUint64, err := strconv.ParseUint(userParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "User ID must be a valid number",
		})
		return
	}
	userId := uint(userIdUint64)

	// 2 Get current user info from JWT
	currentUserID := c.GetUint("userId")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User ID not found in JWT",
		})
		return
	}
	role := c.GetString("role")

	// 3 Permission check
	if role != "admin" && userId != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only view your own support tickets",
		})
		return
	}

	// 4 Fetch tickets
	tickets, err := sh.SupportService.GetSupportTicketsByUserId(userId)
	if err != nil {
		log.Printf("Error fetching tickets for user %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve support tickets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tickets,
		"count": len(tickets),
	})
}

// GetOpenTickets retrieves all open support tickets
func (sh *SupportHandler) GetOpenTickets(c *gin.Context) {
	tickets, err := sh.SupportService.GetOpenTickets()
	if err != nil {
		log.Printf("Error fetching open tickets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve open support tickets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tickets,
		"count": len(tickets),
	})
}

// UpdateTicket updates a support ticket
func (sh *SupportHandler) UpdateTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Ticket ID must be a valid number",
		})
		return
	}

	ticket, err := sh.SupportService.GetSupportTicketById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Support ticket not found",
			})
			return
		}
		log.Printf("Error fetching ticket %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve support ticket",
		})
		return
	}

	userID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only update your own tickets",
		})
		return
	}

	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	if req.Subject != "" {
		ticket.Subject = req.Subject
	}
	if req.Message != "" {
		ticket.Message = req.Message
	}
	if req.Status != "" {
		// Only admins can change status to anything other than "open"
		if role != "admin" && req.Status != "open" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Only admins can change ticket status",
			})
			return
		}
		ticket.Status = req.Status
	}

	if err := sh.SupportService.UpdateSupportTicket(ticket); err != nil {
		log.Printf("Error updating ticket %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to update support ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket updated successfully",
		"data":    ticket,
	})
}

// DeleteTicket deletes a support ticket
func (sh *SupportHandler) DeleteTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Ticket ID must be a valid number",
		})
		return
	}

	ticket, err := sh.SupportService.GetSupportTicketById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Support ticket not found",
			})
			return
		}
		log.Printf("Error fetching ticket %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve support ticket",
		})
		return
	}

	userID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && ticket.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only delete your own tickets",
		})
		return
	}

	if err := sh.SupportService.DeleteSupportTicket(uint(id)); err != nil {
		log.Printf("Error deleting ticket %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to delete support ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

// GetAllTickets retrieves all support tickets with pagination (admin only)
func (sh *SupportHandler) GetAllTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	tickets, err := sh.SupportService.GetAllTickets()
	if err != nil {
		log.Printf("Error fetching all tickets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve support tickets",
		})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(tickets) {
		start = len(tickets)
	}
	if end > len(tickets) {
		end = len(tickets)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tickets[start:end],
		"page":  page,
		"limit": limit,
		"total": len(tickets),
	})
}
