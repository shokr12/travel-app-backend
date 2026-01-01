// repos/visa_repo.go
package repos

import (
	"Visa/models"

	"gorm.io/gorm"
)

type VisaRepo struct {
	db *gorm.DB
}

func NewVisaRepo(db *gorm.DB) *VisaRepo {
	return &VisaRepo{db: db}
}

// GetAllVisas retrieves all visa applications from the database
func (vr *VisaRepo) GetAllVisas() ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// GetVisaById retrieves a visa application by its ID
func (vr *VisaRepo) GetVisaById(id uint) (*models.VisaApplication, error) {
	var visa models.VisaApplication
	if err := vr.db.First(&visa, id).Error; err != nil {
		return nil, err
	}
	return &visa, nil
}

// CreateVisa creates a new visa application in the database
func (vr *VisaRepo) CreateVisa(visa *models.VisaApplication) error {
	return vr.db.Create(visa).Error
}

// UpdateVisa updates an existing visa application
func (vr *VisaRepo) UpdateVisa(visa *models.VisaApplication) error {
	return vr.db.Save(visa).Error
}

// DeleteVisa deletes a visa application by its ID
func (vr *VisaRepo) DeleteVisa(id uint) error {
	return vr.db.Delete(&models.VisaApplication{}, id).Error
}

// FindVisaByCountry retrieves all visa applications for a specific country
func (vr *VisaRepo) FindVisaByCountry(country string) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Where("country = ?", country).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// FindVisaByDestination retrieves all visa applications for a specific destination
func (vr *VisaRepo) FindVisaByDestination(destination string) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Where("destination = ?", destination).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// FindVisaByNationality retrieves all visa applications for a specific nationality
func (vr *VisaRepo) FindVisaByNationality(nationality string) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Where("nationality = ?", nationality).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// GetVisaByUserId retrieves all visa applications for a specific user
func (vr *VisaRepo) GetVisaByUserId(userId uint) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Where("userId = ?", userId).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// GetVisasByStatus retrieves all visa applications with a specific status
// Valid statuses: "pending", "approved", "rejected"
func (vr *VisaRepo) GetVisasByStatus(status string) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Where("status = ?", status).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// GetPendingVisas retrieves all pending visa applications
func (vr *VisaRepo) GetPendingVisas() ([]models.VisaApplication, error) {
	return vr.GetVisasByStatus("pending")
}

// GetApprovedVisas retrieves all approved visa applications
func (vr *VisaRepo) GetApprovedVisas() ([]models.VisaApplication, error) {
	return vr.GetVisasByStatus("approved")
}

// GetRejectedVisas retrieves all rejected visa applications
func (vr *VisaRepo) GetRejectedVisas() ([]models.VisaApplication, error) {
	return vr.GetVisasByStatus("rejected")
}

// FindVisaByVisaType retrieves all visa applications for a specific visa type
func (vr *VisaRepo) FindVisaByVisaType(visaType string) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Where("visa_type = ?", visaType).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}

// GetRecentVisas retrieves the most recent visa applications (ordered by creation date)
func (vr *VisaRepo) GetRecentVisas(limit int) ([]models.VisaApplication, error) {
	var visas []models.VisaApplication
	if err := vr.db.Order("created_at desc").Limit(limit).Find(&visas).Error; err != nil {
		return nil, err
	}
	return visas, nil
}
