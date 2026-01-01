// services/flight_service.go
package services

import (
	"Visa/internal/repos"
	"Visa/models"
	"errors"
	"fmt"
	"strconv"
)

type FlightService struct {
	Repo            *repos.FlightRepo
	ReservationRepo *repos.ReservationRepo
}

func NewFlightService(flightRepo *repos.FlightRepo, reservationRepo *repos.ReservationRepo) *FlightService {
	return &FlightService{
		Repo:            flightRepo,
		ReservationRepo: reservationRepo,
	}
}

// BookFlight books a flight for a user
func (fs *FlightService) BookFlight(userId uint, flightId uint) error {
	// Validate input
	if userId == 0 {
		return errors.New("invalid user ID")
	}
	if flightId == 0 {
		return errors.New("invalid flight ID")
	}

	// Get flight details
	flight, err := fs.Repo.GetFlightById(flightId)
	if err != nil {
		return fmt.Errorf("flight not found: %w", err)
	}

	// Check seat availability
	if flight.SeatsAvailable <= 0 {
		return errors.New("no seats available for this flight")
	}

	// Check if user already has an active booking for this flight
	userFlights, err := fs.Repo.FindByUser(userId)
	if err == nil {
		for _, f := range userFlights {
			if f.ID == flightId {
				return errors.New("you already have a booking for this flight")
			}
		}
	}

	// Create reservation
	flightIDStr := strconv.FormatUint(uint64(flightId), 10)
	res := &models.Reservation{
		UserID:   strconv.FormatUint(uint64(userId), 10),
		FlightID: &flightIDStr,
		Status:   "booked",
	}

	if err := fs.ReservationRepo.CreateReservation(res); err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// Reduce available seats
	flight.SeatsAvailable -= 1
	if err := fs.Repo.UpdateFlight(flight); err != nil {
		// Try to rollback reservation if flight update fails
		fs.ReservationRepo.DeleteReservation(res.ID)
		return fmt.Errorf("failed to update flight availability: %w", err)
	}

	return nil
}

// CancelFlight cancels a flight booking for a user
func (fs *FlightService) CancelFlight(userId uint, flightId uint) error {
	// Validate input
	if userId == 0 {
		return errors.New("invalid user ID")
	}
	if flightId == 0 {
		return errors.New("invalid flight ID")
	}

	// Get flight details
	flight, err := fs.Repo.GetFlightById(flightId)
	if err != nil {
		return fmt.Errorf("flight not found: %w", err)
	}

	// Verify user has a booking for this flight
	userFlights, err := fs.Repo.FindByUser(userId)
	if err != nil {
		return fmt.Errorf("failed to verify booking: %w", err)
	}

	hasBooking := false
	for _, f := range userFlights {
		if f.ID == flightId {
			hasBooking = true
			break
		}
	}

	if !hasBooking {
		return errors.New("no active booking found for this flight")
	}

	// Create cancellation reservation record
	flightIDStr := strconv.FormatUint(uint64(flightId), 10)
	res := &models.Reservation{
		UserID:   strconv.FormatUint(uint64(userId), 10),
		FlightID: &flightIDStr,
		Status:   "cancelled",
	}

	if err := fs.ReservationRepo.CreateReservation(res); err != nil {
		return fmt.Errorf("failed to create cancellation record: %w", err)
	}

	// Increase available seats
	flight.SeatsAvailable += 1
	if err := fs.Repo.UpdateFlight(flight); err != nil {
		return fmt.Errorf("failed to update flight availability: %w", err)
	}

	return nil
}

// GetFlights retrieves all flights
func (fs *FlightService) GetFlights() ([]models.Flight, error) {
	flights, err := fs.Repo.GetAllFlights()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve flights: %w", err)
	}
	return flights, nil
}

// GetFlightById retrieves a flight by its ID
func (fs *FlightService) GetFlightById(id uint) (*models.Flight, error) {
	if id == 0 {
		return nil, errors.New("invalid flight ID")
	}

	flight, err := fs.Repo.GetFlightById(id)
	if err != nil {
		return nil, fmt.Errorf("flight not found: %w", err)
	}
	return flight, nil
}

// GetFlightsByCity retrieves all flights for a specific city
func (fs *FlightService) GetFlightsByCity(city string) ([]models.Flight, error) {
	if city == "" {
		return nil, errors.New("city is required")
	}

	flights, err := fs.Repo.FindFlightByCity(city)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve flights for city %s: %w", city, err)
	}
	return flights, nil
}

// GetFlightsByDepartDate retrieves all flights for a specific departure date
func (fs *FlightService) GetFlightsByDepartDate(date string) ([]models.Flight, error) {
	if date == "" {
		return nil, errors.New("departure date is required")
	}

	flights, err := fs.Repo.FindByDepartDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve flights for date %s: %w", date, err)
	}
	return flights, nil
}

// GetFlightsByArriveDate retrieves all flights for a specific arrival date
func (fs *FlightService) GetFlightsByArriveDate(date string) ([]models.Flight, error) {
	if date == "" {
		return nil, errors.New("arrival date is required")
	}

	flights, err := fs.Repo.FindByArriveDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve flights for date %s: %w", date, err)
	}
	return flights, nil
}

// GetFlightsByClass retrieves all flights for a specific class
func (fs *FlightService) GetFlightsByClass(class string) ([]models.Flight, error) {
	if class == "" {
		return nil, errors.New("flight class is required")
	}

	// Validate class
	validClasses := map[string]bool{
		"economy":  true,
		"business": true,
		"first":    true,
	}
	if !validClasses[class] {
		return nil, errors.New("invalid flight class. Must be economy, business, or first")
	}

	flights, err := fs.Repo.FindByClass(class)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve flights for class %s: %w", class, err)
	}
	return flights, nil
}

// GetDirectFlights retrieves direct or connecting flights
func (fs *FlightService) GetDirectFlights(direct bool) ([]models.Flight, error) {
	flights, err := fs.Repo.DirectFlight(direct)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve direct flights: %w", err)
	}
	return flights, nil
}

// GetFlightsSortedByPrice retrieves flights sorted by price
func (fs *FlightService) GetFlightsSortedByPrice(ascending bool) ([]models.Flight, error) {
	var flights []models.Flight
	var err error

	if ascending {
		flights, err = fs.Repo.SortFromLowerToUpper()
	} else {
		flights, err = fs.Repo.SortFromHigherToLower()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sorted flights: %w", err)
	}
	return flights, nil
}

// GetFlightsByUser retrieves all flights booked by a specific user
func (fs *FlightService) GetFlightsByUser(userId uint) ([]models.Flight, error) {
	if userId == 0 {
		return nil, errors.New("invalid user ID")
	}

	flights, err := fs.Repo.FindByUser(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user flights: %w", err)
	}
	return flights, nil
}

// CreateFlight creates a new flight (admin only)
func (fs *FlightService) CreateFlight(flight *models.Flight) error {
	if flight == nil {
		return errors.New("flight data is required")
	}

	// Validate flight data
	if flight.City == "" {
		return errors.New("flight city is required")
	}
	if flight.Price <= 0 {
		return errors.New("flight price must be greater than 0")
	}
	if flight.SeatsAvailable < 0 {
		return errors.New("seats available cannot be negative")
	}

	if err := fs.Repo.CreateFlight(flight); err != nil {
		return fmt.Errorf("failed to create flight: %w", err)
	}
	return nil
}

// UpdateFlight updates an existing flight (admin only)
func (fs *FlightService) UpdateFlight(flight *models.Flight) error {
	if flight == nil {
		return errors.New("flight data is required")
	}
	if flight.ID == 0 {
		return errors.New("flight ID is required")
	}

	// Verify flight exists
	_, err := fs.Repo.GetFlightById(flight.ID)
	if err != nil {
		return fmt.Errorf("flight not found: %w", err)
	}

	if err := fs.Repo.UpdateFlight(flight); err != nil {
		return fmt.Errorf("failed to update flight: %w", err)
	}
	return nil
}

// DeleteFlight deletes a flight (admin only)
func (fs *FlightService) DeleteFlight(id uint) error {
	if id == 0 {
		return errors.New("invalid flight ID")
	}

	// Check if flight has active bookings
	// You might want to prevent deletion if there are active bookings
	if err := fs.Repo.DeleteFlight(id); err != nil {
		return fmt.Errorf("failed to delete flight: %w", err)
	}
	return nil
}
