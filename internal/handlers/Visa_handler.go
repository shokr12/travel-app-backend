// handlers/visa_handler.go
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

type VisaHandler struct {
	VisaService *services.VisaService
}

func NewVisaHandler(visaService *services.VisaService) *VisaHandler {
	return &VisaHandler{VisaService: visaService}
}

type CreateVisaRequest struct {
	UserID         uint   `json:"userId" binding:"required"`
	VisaType       string `json:"visa_type" binding:"required,min=2,max=50"`
	Destination    string `json:"destination" binding:"required,min=2,max=100"`
	TravelDate     string `json:"travel_date" binding:"required"`
	PassportNumber string `json:"passport_number" binding:"required,min=6,max=20"`
	Nationality    string `json:"nationality" binding:"required,min=2,max=50"`
}

type UpdateVisaRequest struct {
	VisaType       string `json:"visa_type" binding:"omitempty,min=2,max=50"`
	Destination    string `json:"destination" binding:"omitempty,min=2,max=100"`
	TravelDate     string `json:"travel_date" binding:"omitempty"`
	PassportNumber string `json:"passport_number" binding:"omitempty,min=6,max=20"`
	Nationality    string `json:"nationality" binding:"omitempty,min=2,max=50"`
}

// CreateVisa creates a new visa application
func (vh *VisaHandler) CreateVisa(c *gin.Context) {
	var req CreateVisaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	userID := c.GetUint("userId")
	if req.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only create visa applications for yourself",
		})
		return
	}

	visa := models.VisaApplication{
		UserID:         req.UserID,
		VisaType:       req.VisaType,
		Destination:    req.Destination,
		TravelDate:     req.TravelDate,
		PassportNumber: req.PassportNumber,
		Nationality:    req.Nationality,
		Status:         "pending",
	}

	if err := vh.VisaService.CreateVisa(&visa); err != nil {
		log.Printf("Error creating visa: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to create visa application",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Visa application submitted successfully",
		"data":    visa,
	})
}

// GetVisaById retrieves a visa application by ID
func (vh *VisaHandler) GetVisaById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Visa ID must be a valid number",
		})
		return
	}

	visa, err := vh.VisaService.GetVisaById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Visa application not found",
			})
			return
		}
		log.Printf("Error fetching visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve visa application",
		})
		return
	}

	userID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && visa.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You don't have permission to view this visa",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": visa})
}

// GetVisasByUser retrieves all visa applications for a specific user
func (vh *VisaHandler) GetVisasByUser(c *gin.Context) {
	userParam := c.Param("userId")
	userId, err := strconv.ParseUint(userParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "User ID must be a valid number",
		})
		return
	}

	currentUserID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && uint(userId) != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only view your own visa applications",
		})
		return
	}

	visas, err := vh.VisaService.GetVisaByUserId(uint(userId))
	if err != nil {
		log.Printf("Error fetching visas for user %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve visa applications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  visas,
		"count": len(visas),
	})
}

// UpdateVisa updates an existing visa application
func (vh *VisaHandler) UpdateVisa(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Visa ID must be a valid number",
		})
		return
	}

	existingVisa, err := vh.VisaService.GetVisaById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Visa application not found",
			})
			return
		}
		log.Printf("Error fetching visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve visa application",
		})
		return
	}

	userID := c.GetUint("userId")
	if existingVisa.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only update your own visa applications",
		})
		return
	}

	if existingVisa.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_status",
			"message": "Only pending visa applications can be updated",
		})
		return
	}

	var req UpdateVisaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	if req.VisaType != "" {
		existingVisa.VisaType = req.VisaType
	}
	if req.Destination != "" {
		existingVisa.Destination = req.Destination
	}
	if req.TravelDate != "" {
		existingVisa.TravelDate = req.TravelDate
	}
	if req.PassportNumber != "" {
		existingVisa.PassportNumber = req.PassportNumber
	}
	if req.Nationality != "" {
		existingVisa.Nationality = req.Nationality
	}

	if err := vh.VisaService.UpdateVisa(existingVisa); err != nil {
		log.Printf("Error updating visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to update visa application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Visa application updated successfully",
		"data":    existingVisa,
	})
}

// DeleteVisa deletes a visa application
func (vh *VisaHandler) DeleteVisa(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Visa ID must be a valid number",
		})
		return
	}

	visa, err := vh.VisaService.GetVisaById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Visa application not found",
			})
			return
		}
		log.Printf("Error fetching visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve visa application",
		})
		return
	}

	userID := c.GetUint("userId")
	role := c.GetString("role")
	if role != "admin" && visa.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "forbidden",
			"message": "You can only delete your own visa applications",
		})
		return
	}

	if err := vh.VisaService.DeleteVisa(uint(id)); err != nil {
		log.Printf("Error deleting visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to delete visa application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visa application deleted successfully"})
}

// GetAllVisa retrieves all visa applications with pagination (admin only)
func (vh *VisaHandler) GetAllVisa(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	visas, err := vh.VisaService.GetAllVisa()
	if err != nil {
		log.Printf("Error fetching all visas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve visa applications",
		})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(visas) {
		start = len(visas)
	}
	if end > len(visas) {
		end = len(visas)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  visas[start:end],
		"page":  page,
		"limit": limit,
		"total": len(visas),
	})
}

// GetApprovedVisas retrieves all approved visa applications (admin only)
func (vh *VisaHandler) GetApprovedVisas(c *gin.Context) {
	visas, err := vh.VisaService.GetApprovedVisas()
	if err != nil {
		log.Printf("Error fetching approved visas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve approved visas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  visas,
		"count": len(visas),
	})
}

// GetPendingVisas retrieves all pending visa applications (admin only)
func (vh *VisaHandler) GetPendingVisas(c *gin.Context) {
	visas, err := vh.VisaService.GetPendingVisas()
	if err != nil {
		log.Printf("Error fetching pending visas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve pending visas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  visas,
		"count": len(visas),
	})
}

// GetRejectedVisas retrieves all rejected visa applications (admin only)
func (vh *VisaHandler) GetRejectedVisas(c *gin.Context) {
	visas, err := vh.VisaService.GetRejectedVisas()
	if err != nil {
		log.Printf("Error fetching rejected visas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve rejected visas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  visas,
		"count": len(visas),
	})
}

// ApproveVisa approves a visa application (admin only)
func (vh *VisaHandler) ApproveVisa(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Visa ID must be a valid number",
		})
		return
	}

	if err := vh.VisaService.ApproveVisa(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Visa application not found",
			})
			return
		}
		log.Printf("Error approving visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to approve visa application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visa approved successfully"})
}

// RejectVisa rejects a visa application (admin only)
func (vh *VisaHandler) RejectVisa(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Visa ID must be a valid number",
		})
		return
	}

	if err := vh.VisaService.RejectVisa(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Visa application not found",
			})
			return
		}
		log.Printf("Error rejecting visa %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to reject visa application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visa rejected successfully"})
}
