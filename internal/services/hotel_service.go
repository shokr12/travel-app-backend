// services/hotel_service.go
package services

import (
	"Visa/internal/repos"
	"Visa/models"
	"errors"
	"fmt"
	"strconv"
)

type HotelService struct {
	Repo            *repos.HotelRepo
	ReservationRepo *repos.ReservationRepo
}

func NewHotelService(hotelRepo *repos.HotelRepo, reservationRepo *repos.ReservationRepo) *HotelService {
	return &HotelService{
		Repo:            hotelRepo,
		ReservationRepo: reservationRepo,
	}
}

// CreateHotel creates a new hotel (admin only)
func (hs *HotelService) CreateHotel(hotel *models.Hotel) error {
	if hotel == nil {
		return errors.New("hotel data is required")
	}

	// Validate hotel data
	if hotel.Name == "" {
		return errors.New("hotel name is required")
	}
	if hotel.City == "" {
		return errors.New("hotel city is required")
	}
	if hotel.Address == "" {
		return errors.New("hotel address is required")
	}
	if hotel.PricePerNight <= 0 {
		return errors.New("price per night must be greater than 0")
	}
	if hotel.AvailableRooms < 0 {
		return errors.New("available rooms cannot be negative")
	}

	if err := hs.Repo.CreateHotel(hotel); err != nil {
		return fmt.Errorf("failed to create hotel: %w", err)
	}
	return nil
}

// UpdateHotel updates an existing hotel (admin only)
func (hs *HotelService) UpdateHotel(hotel *models.Hotel) error {
	if hotel == nil {
		return errors.New("hotel data is required")
	}
	if hotel.ID == 0 {
		return errors.New("hotel ID is required")
	}

	// Verify hotel exists
	_, err := hs.Repo.GetHotelById(hotel.ID)
	if err != nil {
		return fmt.Errorf("hotel not found: %w", err)
	}

	if err := hs.Repo.UpdateHotel(hotel); err != nil {
		return fmt.Errorf("failed to update hotel: %w", err)
	}
	return nil
}

// DeleteHotel deletes a hotel (admin only)
func (hs *HotelService) DeleteHotel(hotelId uint) error {
	if hotelId == 0 {
		return errors.New("invalid hotel ID")
	}

	// Verify hotel exists
	_, err := hs.Repo.GetHotelById(hotelId)
	if err != nil {
		return fmt.Errorf("hotel not found: %w", err)
	}

	if err := hs.Repo.DeleteHotel(hotelId); err != nil {
		return fmt.Errorf("failed to delete hotel: %w", err)
	}
	return nil
}

// BookHotel books a hotel room for a user
func (hs *HotelService) BookHotel(userId uint, hotelId uint) error {
	// Validate input
	if userId == 0 {
		return errors.New("invalid user ID")
	}
	if hotelId == 0 {
		return errors.New("invalid hotel ID")
	}

	// Get hotel details
	hotel, err := hs.Repo.GetHotelById(hotelId)
	if err != nil {
		return fmt.Errorf("hotel not found: %w", err)
	}

	// Check room availability
	if hotel.AvailableRooms <= 0 {
		return errors.New("no rooms available at this hotel")
	}

	// Check if user already has an active booking for this hotel
	userHotels, err := hs.Repo.GetHotelsByUser(userId)
	if err == nil {
		for _, h := range userHotels {
			if h.ID == hotelId {
				return errors.New("you already have a booking at this hotel")
			}
		}
	}

	// Create reservation
	hotelIDStr := strconv.FormatUint(uint64(hotelId), 10)
	res := &models.Reservation{
		UserID:  strconv.FormatUint(uint64(userId), 10),
		HotelID: hotelIDStr,
		Status:  "booked",
	}

	if err := hs.ReservationRepo.CreateReservation(res); err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// Reduce available rooms
	hotel.AvailableRooms -= 1
	if err := hs.Repo.UpdateHotel(hotel); err != nil {
		// Try to rollback reservation if hotel update fails
		hs.ReservationRepo.DeleteReservation(res.ID)
		return fmt.Errorf("failed to update hotel availability: %w", err)
	}

	return nil
}

// CancelHotel cancels a hotel booking for a user
func (hs *HotelService) CancelHotel(userId uint, hotelId uint) error {
	// Validate input
	if userId == 0 {
		return errors.New("invalid user ID")
	}
	if hotelId == 0 {
		return errors.New("invalid hotel ID")
	}

	// Get hotel details
	hotel, err := hs.Repo.GetHotelById(hotelId)
	if err != nil {
		return fmt.Errorf("hotel not found: %w", err)
	}

	// Verify user has a booking for this hotel
	userHotels, err := hs.Repo.GetHotelsByUser(userId)
	if err != nil {
		return fmt.Errorf("failed to verify booking: %w", err)
	}

	hasBooking := false
	for _, h := range userHotels {
		if h.ID == hotelId {
			hasBooking = true
			break
		}
	}

	if !hasBooking {
		return errors.New("no active booking found for this hotel")
	}

	// Check if hotel allows free cancellation
	if !hotel.FreeCancellation {
		return errors.New("this hotel does not allow free cancellation")
	}

	// Create cancellation record
	hotelIDStr := strconv.FormatUint(uint64(hotelId), 10)
	res := &models.Reservation{
		UserID:  strconv.FormatUint(uint64(userId), 10),
		HotelID: hotelIDStr,
		Status:  "cancelled",
	}

	if err := hs.ReservationRepo.CreateReservation(res); err != nil {
		return fmt.Errorf("failed to create cancellation record: %w", err)
	}

	// Increase available rooms
	hotel.AvailableRooms += 1
	if err := hs.Repo.UpdateHotel(hotel); err != nil {
		return fmt.Errorf("failed to update hotel availability: %w", err)
	}

	return nil
}

// GetAllHotels retrieves all hotels
func (hs *HotelService) GetAllHotels() ([]models.Hotel, error) {
	hotels, err := hs.Repo.GetAllHotels()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hotels: %w", err)
	}
	return hotels, nil
}

// GetHotelById retrieves a hotel by its ID
func (hs *HotelService) GetHotelById(hotelId uint) (*models.Hotel, error) {
	if hotelId == 0 {
		return nil, errors.New("invalid hotel ID")
	}

	hotel, err := hs.Repo.GetHotelById(hotelId)
	if err != nil {
		return nil, fmt.Errorf("hotel not found: %w", err)
	}
	return hotel, nil
}

// GetHotelsByCity retrieves all hotels in a specific city
func (hs *HotelService) GetHotelsByCity(city string) ([]models.Hotel, error) {
	if city == "" {
		return nil, errors.New("city is required")
	}

	hotels, err := hs.Repo.FindHotelByCity(city)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hotels for city %s: %w", city, err)
	}
	return hotels, nil
}

// GetHotelsByCheckInDate retrieves hotels by check-in date
func (hs *HotelService) GetHotelsByCheckInDate(date string) ([]models.Hotel, error) {
	if date == "" {
		return nil, errors.New("check-in date is required")
	}

	hotels, err := hs.Repo.FindByCheckInDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hotels for check-in date %s: %w", date, err)
	}
	return hotels, nil
}

// GetHotelsByCheckOutDate retrieves hotels by check-out date
func (hs *HotelService) GetHotelsByCheckOutDate(date string) ([]models.Hotel, error) {
	if date == "" {
		return nil, errors.New("check-out date is required")
	}

	hotels, err := hs.Repo.FindByCheckOutDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hotels for check-out date %s: %w", date, err)
	}
	return hotels, nil
}

// GetHotelOptions retrieves hotels based on free cancellation option
func (hs *HotelService) GetHotelOptions(freeCancellation bool) ([]models.Hotel, error) {
	hotels, err := hs.Repo.HotelOptions(freeCancellation)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hotels with free cancellation option: %w", err)
	}
	return hotels, nil
}

// GetHotelsSortedByPrice retrieves hotels sorted by price
func (hs *HotelService) GetHotelsSortedByPrice(ascending bool) ([]models.Hotel, error) {
	var hotels []models.Hotel
	var err error

	if ascending {
		hotels, err = hs.Repo.SortFromLowerToUpper()
	} else {
		hotels, err = hs.Repo.SortFromHigherToLower()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sorted hotels: %w", err)
	}
	return hotels, nil
}

// GetHotelsByUser retrieves all hotels booked by a specific user
func (hs *HotelService) GetHotelsByUser(userId uint) ([]models.Hotel, error) {
	if userId == 0 {
		return nil, errors.New("invalid user ID")
	}

	hotels, err := hs.Repo.GetHotelsByUser(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user hotels: %w", err)
	}
	return hotels, nil
}

// SearchHotels searches hotels with multiple filters
func (hs *HotelService) SearchHotels(city string, minPrice, maxPrice float64, freeCancellation *bool) ([]models.Hotel, error) {
	hotels, err := hs.GetAllHotels()
	if err != nil {
		return nil, err
	}

	var filtered []models.Hotel
	for _, hotel := range hotels {
		// Filter by city
		if city != "" && hotel.City != city {
			continue
		}

		// Filter by price range
		if minPrice > 0 && hotel.PricePerNight < minPrice {
			continue
		}
		if maxPrice > 0 && hotel.PricePerNight > maxPrice {
			continue
		}

		// Filter by free cancellation
		if freeCancellation != nil && hotel.FreeCancellation != *freeCancellation {
			continue
		}

		filtered = append(filtered, hotel)
	}

	return filtered, nil
}
