// handlers/flight_handler.go
package handlers

import (
	"Visa/internal/services"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FlightHandler struct {
	FlightService *services.FlightService
}

func NewFlightHandler(flightService *services.FlightService) *FlightHandler {
	return &FlightHandler{FlightService: flightService}
}

type BookFlightRequest struct {
	UserID   uint `json:"userId" binding:"required"`
	FlightID uint `json:"flight_id" binding:"required"`
}

// BookFlight books a flight for a user
func (fh *FlightHandler) BookFlight(c *gin.Context) {
	var req BookFlightRequest
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
			"message": "You can only book flights for yourself",
		})
		return
	}

	if err := fh.FlightService.BookFlight(req.UserID, req.FlightID); err != nil {
		log.Printf("Error booking flight %d for user %d: %v", req.FlightID, req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to book flight",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight booked successfully"})
}

// CancelFlight cancels a flight booking
func (fh *FlightHandler) CancelFlight(c *gin.Context) {
	var req BookFlightRequest
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
			"message": "You can only cancel your own flight bookings",
		})
		return
	}

	if err := fh.FlightService.CancelFlight(req.UserID, req.FlightID); err != nil {
		log.Printf("Error cancelling flight %d for user %d: %v", req.FlightID, req.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to cancel flight",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flight cancelled successfully"})
}

// GetFlights retrieves all flights with pagination
func (fh *FlightHandler) GetFlights(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	flights, err := fh.FlightService.GetFlights()
	if err != nil {
		log.Printf("Error fetching flights: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve flights",
		})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if start > len(flights) {
		start = len(flights)
	}
	if end > len(flights) {
		end = len(flights)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  flights[start:end],
		"page":  page,
		"limit": limit,
		"total": len(flights),
	})
}

// GetFlightById retrieves a flight by ID
func (fh *FlightHandler) GetFlightById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "Flight ID must be a valid number",
		})
		return
	}

	flight, err := fh.FlightService.GetFlightById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "not_found",
				"message": "Flight not found",
			})
			return
		}
		log.Printf("Error fetching flight %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve flight",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": flight})
}

// GetFlightsByCity retrieves flights by city
func (fh *FlightHandler) GetFlightsByCity(c *gin.Context) {
	city := c.Param("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "City parameter is required",
		})
		return
	}

	flights, err := fh.FlightService.GetFlightsByCity(city)
	if err != nil {
		log.Printf("Error fetching flights by city %s: %v", city, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve flights",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  flights,
		"count": len(flights),
	})
}

// GetFlightsByDepartDate retrieves flights by departure date
func (fh *FlightHandler) GetFlightsByDepartDate(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_input",
			"message": "Date parameter is required",
		})
		return
	}

	flights, err := fh.FlightService.GetFlightsByDepartDate(date)
	if err != nil {
		log.Printf("Error fetching flights by date %s: %v", date, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve flights",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  flights,
		"count": len(flights),
	})
}

// GetFlightsByUser retrieves all flights booked by a user
func (fh *FlightHandler) GetFlightsByUser(c *gin.Context) {
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
			"message": "You can only view your own flight bookings",
		})
		return
	}

	flights, err := fh.FlightService.GetFlightsByUser(uint(userId))
	if err != nil {
		log.Printf("Error fetching flights for user %d: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "server_error",
			"message": "Unable to retrieve flight bookings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  flights,
		"count": len(flights),
	})
}
