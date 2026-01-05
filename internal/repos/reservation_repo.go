// repos/reservation_repo.go
package repos

import (
	"Visa/models"

	"gorm.io/gorm"
)

type ReservationRepo struct {
	db *gorm.DB
}

func NewReservationRepo(db *gorm.DB) *ReservationRepo {
	return &ReservationRepo{db: db}
}

// GetAllReservations retrieves all reservations from the database
func (rr *ReservationRepo) GetAllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := rr.db.Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationById retrieves a reservation by its ID
func (rr *ReservationRepo) GetReservationById(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := rr.db.First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

// CreateReservation creates a new reservation in the database
func (rr *ReservationRepo) CreateReservation(reservation *models.Reservation) error {
	return rr.db.Create(reservation).Error
}

// UpdateReservation updates an existing reservation
func (rr *ReservationRepo) UpdateReservation(reservation *models.Reservation) error {
	return rr.db.Save(reservation).Error
}

// DeleteReservation deletes a reservation by its ID
func (rr *ReservationRepo) DeleteReservation(id uint) error {
	return rr.db.Delete(&models.Reservation{}, id).Error
}

// GetReservationsByUserId retrieves all reservations for a specific user
func (rr *ReservationRepo) GetReservationsByUserId(userId uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := rr.db.Where("user_id = ?", userId).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByHotelId retrieves all reservations for a specific hotel
func (rr *ReservationRepo) GetReservationsByHotelId(hotelId uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := rr.db.Where("hotel_id = ?", hotelId).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByFlightId retrieves all reservations for a specific flight
func (rr *ReservationRepo) GetReservationsByFlightId(flightId uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := rr.db.Where("flight_id = ?", flightId).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByStatus retrieves all reservations with a specific status
func (rr *ReservationRepo) GetReservationsByStatus(status string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	if err := rr.db.Where("status = ?", status).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}
