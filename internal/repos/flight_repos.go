// repos/flight_repo.go
package repos

import (
	"Visa/models"

	"gorm.io/gorm"
)

type FlightRepo struct {
	db *gorm.DB
}

func NewFlightRepo(db *gorm.DB) *FlightRepo {
	return &FlightRepo{db: db}
}

// GetAllFlights retrieves all flights from the database
func (fr *FlightRepo) GetAllFlights() ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// GetFlightById retrieves a flight by its ID
func (fr *FlightRepo) GetFlightById(id uint) (*models.Flight, error) {
	var flight models.Flight
	if err := fr.db.First(&flight, id).Error; err != nil {
		return nil, err
	}
	return &flight, nil
}

// CreateFlight creates a new flight in the database
func (fr *FlightRepo) CreateFlight(flight *models.Flight) error {
	return fr.db.Create(flight).Error
}

// UpdateFlight updates an existing flight
func (fr *FlightRepo) UpdateFlight(flight *models.Flight) error {
	return fr.db.Save(flight).Error
}

// DeleteFlight deletes a flight by its ID
func (fr *FlightRepo) DeleteFlight(id uint) error {
	return fr.db.Delete(&models.Flight{}, id).Error
}

// FindFlightByCity retrieves all flights for a specific city
func (fr *FlightRepo) FindFlightByCity(city string) ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Where("city = ?", city).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// FindByDepartDate retrieves all flights for a specific departure date
func (fr *FlightRepo) FindByDepartDate(date string) ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Where("depart_date = ?", date).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// FindByArriveDate retrieves all flights for a specific arrival date
func (fr *FlightRepo) FindByArriveDate(date string) ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Where("arrive_date = ?", date).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// FindByClass retrieves all flights for a specific class (economy, business, first)
func (fr *FlightRepo) FindByClass(class string) ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Where("class = ?", class).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// DirectFlight retrieves all direct or connecting flights based on the flag
func (fr *FlightRepo) DirectFlight(checked bool) ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Where("direct_flight = ?", checked).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// SortFromHigherToLower retrieves all flights sorted by price (highest to lowest)
func (fr *FlightRepo) SortFromHigherToLower() ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Order("price desc").Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// SortFromLowerToUpper retrieves all flights sorted by price (lowest to highest)
func (fr *FlightRepo) SortFromLowerToUpper() ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Order("price asc").Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}

// FindByUser retrieves all flights booked by a specific user
func (fr *FlightRepo) FindByUser(userId uint) ([]models.Flight, error) {
	var flights []models.Flight
	if err := fr.db.Where("userId = ?", userId).Find(&flights).Error; err != nil {
		return nil, err
	}
	return flights, nil
}
