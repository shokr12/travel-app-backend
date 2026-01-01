// handlers/hotel_handler.go
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

type HotelHandler struct {
	HotelService *services.HotelService
}

func NewHotelHandler(hotelService *services.HotelService) *HotelHandler {
	return &HotelHandler{HotelService: hotelService}
}

type CreateHotelRequest struct {
	Name             string  `json:"name" binding:"required,min=2,max=200"`
	City             string  `json:"city" binding:"required,min=2,max=100"`
	Address          string  `json:"address" binding:"required,min=5,max=300"`
	PricePerNight    float64 `json:"price_per_night" binding:"required,gt=0"`
	CheckInDate      string  `json:"check_in_date" binding:"required"`
	CheckOutDate     string  `json:"check_out_date" binding:"required"`
	FreeCancellation bool    `json:"free_cancellation"`
	AvailableRooms   int     `json:"available_rooms" binding:"required,gte=0"`
}

type BookHotelRequest struct {
	UserID  uint `json:"userId" binding:"required"`
	HotelID uint `json:"hotel_id" binding:"required"`
}

// CreateHotel creates a new hotel (admin only)
func (hh *HotelHandler) CreateHotel(c *gin.Context) {
	var req CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	hotel := models.Hotel{
		Name:             req.Name,
		City:             req.City,
		Address:          req.Address,
		PricePerNight:    req.PricePerNight,
		CheckInDate:      req.CheckInDate,
		CheckOutDate:     req.CheckOutDate,
		FreeCancellation: req.FreeCancellation,
		AvailableRooms:   req.AvailableRooms,
	}

	if err := hh.HotelService.CreateHotel(&hotel); err != nil {
		log.Printf("Error creating hotel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to create hotel",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Hotel created successfully",
		"data":    hotel,
	})
}

// UpdateHotel updates an existing hotel (admin only)
func (hh *HotelHandler) UpdateHotel(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Hotel ID must be a valid number",
		})
		return
	}

	var hotel models.Hotel
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Please check your input data",
			"details": err.Error(),
		})
		return
	}

	hotel.ID = uint(id)

	if err := hh.HotelService.UpdateHotel(&hotel); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Hotel not found",
			})
			return
		}
		log.Printf("Error updating hotel %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to update hotel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel updated successfully",
		"data":    hotel,
	})
}

// DeleteHotel deletes a hotel (admin only)
func (hh *HotelHandler) DeleteHotel(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Hotel ID must be a valid number",
		})
		return
	}

	if err := hh.HotelService.DeleteHotel(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Hotel not found",
			})
			return
		}
		log.Printf("Error deleting hotel %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to delete hotel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel deleted successfully"})
}

// GetAllHotels retrieves all hotels with pagination
func (hh *HotelHandler) GetAllHotels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	hotels, err := hh.HotelService.GetAllHotels()
	if err != nil {
		log.Printf("Error fetching hotels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotels",
		})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(hotels) {
		start = len(hotels)
	}
	if end > len(hotels) {
		end = len(hotels)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  hotels[start:end],
		"page":  page,
		"limit": limit,
		"total": len(hotels),
	})
}

// GetHotelById retrieves a hotel by ID
func (hh *HotelHandler) GetHotelById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Hotel ID must be a valid number",
		})
		return
	}

	hotel, err := hh.HotelService.GetHotelById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Hotel not found",
			})
			return
		}
		log.Printf("Error fetching hotel %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": hotel})
}

// GetHotelsByCity retrieves hotels by city
func (hh *HotelHandler) GetHotelsByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "City parameter is required",
		})
		return
	}

	hotels, err := hh.HotelService.GetHotelsByCity(city)
	if err != nil {
		log.Printf("Error fetching hotels by city %s: %v", city, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  hotels,
		"count": len(hotels),
	})
}

// GetHotelsByCheckInDate retrieves hotels by check-in date
func (hh *HotelHandler) GetHotelsByCheckInDate(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Date parameter is required",
		})
		return
	}

	hotels, err := hh.HotelService.GetHotelsByCheckInDate(date)
	if err != nil {
		log.Printf("Error fetching hotels by check-in date %s: %v", date, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  hotels,
		"count": len(hotels),
	})
}

// GetHotelsByCheckOutDate retrieves hotels by check-out date
func (hh *HotelHandler) GetHotelsByCheckOutDate(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Date parameter is required",
		})
		return
	}

	hotels, err := hh.HotelService.GetHotelsByCheckOutDate(date)
	if err != nil {
		log.Printf("Error fetching hotels by check-out date %s: %v", date, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  hotels,
		"count": len(hotels),
	})
}

// GetHotelsByFreeCancellation retrieves hotels by free cancellation option
func (hh *HotelHandler) GetHotelsByFreeCancellation(c *gin.Context) {
	freeCancellation := c.DefaultQuery("free_cancellation", "true")
	checked := freeCancellation == "true"

	hotels, err := hh.HotelService.GetHotelOptions(checked)
	if err != nil {
		log.Printf("Error fetching hotels by free cancellation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  hotels,
		"count": len(hotels),
	})
}

// BookHotel books a hotel for a user
func (hh *HotelHandler) BookHotel(c *gin.Context) {
	var req BookHotelRequest
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
			"message": "You can only book hotels for yourself",
		})
		return
	}

	if err := hh.HotelService.BookHotel(req.UserID, req.HotelID); err != nil {
		log.Printf("Error booking hotel %d for user %d: %v", req.HotelID, req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to book hotel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel booked successfully"})
}

// CancelHotel cancels a hotel booking
func (hh *HotelHandler) CancelHotel(c *gin.Context) {
	var req BookHotelRequest
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
			"message": "You can only cancel your own bookings",
		})
		return
	}

	if err := hh.HotelService.CancelHotel(req.UserID, req.HotelID); err != nil {
		log.Printf("Error cancelling hotel %d for user %d: %v", req.HotelID, req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to cancel hotel booking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel booking cancelled successfully"})
}

// GetHotelsByUser retrieves all hotels booked by a user
func (hh *HotelHandler) GetHotelsByUser(c *gin.Context) {
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
			"message": "You can only view your own hotel bookings",
		})
		return
	}

	hotels, err := hh.HotelService.GetHotelsByUser(uint(userId))
	if err != nil {
		log.Printf("Error fetching hotels for user %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve hotel bookings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  hotels,
		"count": len(hotels),
	})
}
