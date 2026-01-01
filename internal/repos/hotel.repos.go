// repos/hotel_repo.go
package repos

import (
	"Visa/models"

	"gorm.io/gorm"
)

type HotelRepo struct {
	db *gorm.DB
}

func NewHotelRepo(db *gorm.DB) *HotelRepo {
	return &HotelRepo{db: db}
}

// GetAllHotels retrieves all hotels from the database
func (hr *HotelRepo) GetAllHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// GetHotelById retrieves a hotel by its ID
func (hr *HotelRepo) GetHotelById(id uint) (*models.Hotel, error) {
	var hotel models.Hotel
	if err := hr.db.First(&hotel, id).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}

// CreateHotel creates a new hotel in the database
func (hr *HotelRepo) CreateHotel(hotel *models.Hotel) error {
	return hr.db.Create(hotel).Error
}

// UpdateHotel updates an existing hotel
func (hr *HotelRepo) UpdateHotel(hotel *models.Hotel) error {
	return hr.db.Save(hotel).Error
}

// DeleteHotel deletes a hotel by its ID
func (hr *HotelRepo) DeleteHotel(id uint) error {
	return hr.db.Delete(&models.Hotel{}, id).Error
}

// FindHotelByCity retrieves all hotels in a specific city
func (hr *HotelRepo) FindHotelByCity(city string) ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Where("city = ?", city).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// FindByCheckInDate retrieves all hotels with a specific check-in date
func (hr *HotelRepo) FindByCheckInDate(date string) ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Where("check_in_date = ?", date).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// FindByCheckOutDate retrieves all hotels with a specific check-out date
func (hr *HotelRepo) FindByCheckOutDate(date string) ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Where("check_out_date = ?", date).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// SortFromHigherToLower retrieves all hotels sorted by price (highest to lowest)
func (hr *HotelRepo) SortFromHigherToLower() ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Order("price desc").Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// SortFromLowerToUpper retrieves all hotels sorted by price (lowest to highest)
func (hr *HotelRepo) SortFromLowerToUpper() ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Order("price asc").Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// HotelOptions retrieves hotels based on free cancellation option
func (hr *HotelRepo) HotelOptions(checked bool) ([]models.Hotel, error) {
	freeCancellation := false
	if checked {
		freeCancellation = true
	}
	var hotels []models.Hotel
	if err := hr.db.Where("free_cancellation = ?", freeCancellation).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}

// GetHotelsByUser retrieves all hotels booked by a specific user
func (hr *HotelRepo) GetHotelsByUser(userId uint) ([]models.Hotel, error) {
	var hotels []models.Hotel
	if err := hr.db.Where("userId = ?", userId).Find(&hotels).Error; err != nil {
		return nil, err
	}
	return hotels, nil
}
