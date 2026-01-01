// services/visa_service.go
package services

import (
	"Visa/internal/repos"
	"Visa/models"
	"errors"
	"fmt"
	"strings"
	"time"
)

type VisaService struct {
	Repo *repos.VisaRepo
}

func NewVisaService(visaRepo *repos.VisaRepo) *VisaService {
	return &VisaService{Repo: visaRepo}
}

// CreateVisa creates a new visa application
func (vs *VisaService) CreateVisa(visa *models.VisaApplication) error {
	if visa == nil {
		return errors.New("visa application data is required")
	}

	// Validate visa data
	if err := vs.validateVisaApplication(visa); err != nil {
		return err
	}

	// Set default status if not provided
	if visa.Status == "" {
		visa.Status = "pending"
	}

	// Validate status
	if visa.Status != "pending" {
		return errors.New("new visa applications must have 'pending' status")
	}

	if err := vs.Repo.CreateVisa(visa); err != nil {
		return fmt.Errorf("failed to create visa application: %w", err)
	}
	return nil
}

// GetVisaById retrieves a visa application by its ID
func (vs *VisaService) GetVisaById(id uint) (*models.VisaApplication, error) {
	if id == 0 {
		return nil, errors.New("invalid visa ID")
	}

	visa, err := vs.Repo.GetVisaById(id)
	if err != nil {
		return nil, fmt.Errorf("visa application not found: %w", err)
	}
	return visa, nil
}

// GetAllVisa retrieves all visa applications
func (vs *VisaService) GetAllVisa() ([]models.VisaApplication, error) {
	visas, err := vs.Repo.GetAllVisas()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve visa applications: %w", err)
	}
	return visas, nil
}

// GetVisaByUserId retrieves all visa applications for a specific user
func (vs *VisaService) GetVisaByUserId(userId uint) ([]models.VisaApplication, error) {
	if userId == 0 {
		return nil, errors.New("invalid user ID")
	}

	visas, err := vs.Repo.GetVisaByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user visa applications: %w", err)
	}
	return visas, nil
}

// UpdateVisa updates an existing visa application
func (vs *VisaService) UpdateVisa(visa *models.VisaApplication) error {
	if visa == nil {
		return errors.New("visa application data is required")
	}
	if visa.ID == 0 {
		return errors.New("visa ID is required")
	}

	// Verify visa exists
	existingVisa, err := vs.Repo.GetVisaById(visa.ID)
	if err != nil {
		return fmt.Errorf("visa application not found: %w", err)
	}

	// Only pending visas can be updated by users
	if existingVisa.Status != "pending" {
		return errors.New("only pending visa applications can be updated")
	}

	// Validate updated fields
	if visa.VisaType != "" && visa.VisaType != existingVisa.VisaType {
		if len(visa.VisaType) < 2 || len(visa.VisaType) > 50 {
			return errors.New("visa type must be between 2 and 50 characters")
		}
		existingVisa.VisaType = visa.VisaType
	}

	if visa.Destination != "" && visa.Destination != existingVisa.Destination {
		if len(visa.Destination) < 2 || len(visa.Destination) > 100 {
			return errors.New("destination must be between 2 and 100 characters")
		}
		existingVisa.Destination = visa.Destination
	}

	if visa.Nationality != "" && visa.Nationality != existingVisa.Nationality {
		if len(visa.Nationality) < 2 || len(visa.Nationality) > 50 {
			return errors.New("nationality must be between 2 and 50 characters")
		}
		existingVisa.Nationality = visa.Nationality
	}

	if visa.PassportNumber != "" && visa.PassportNumber != existingVisa.PassportNumber {
		if len(visa.PassportNumber) < 6 || len(visa.PassportNumber) > 20 {
			return errors.New("passport number must be between 6 and 20 characters")
		}
		existingVisa.PassportNumber = visa.PassportNumber
	}

	if visa.TravelDate != "" && visa.TravelDate != existingVisa.TravelDate {
		existingVisa.TravelDate = visa.TravelDate
	}

	if err := vs.Repo.UpdateVisa(existingVisa); err != nil {
		return fmt.Errorf("failed to update visa application: %w", err)
	}
	return nil
}

// DeleteVisa deletes a visa application
func (vs *VisaService) DeleteVisa(id uint) error {
	if id == 0 {
		return errors.New("invalid visa ID")
	}

	// Verify visa exists
	visa, err := vs.Repo.GetVisaById(id)
	if err != nil {
		return fmt.Errorf("visa application not found: %w", err)
	}

	// Only pending or rejected visas can be deleted
	if visa.Status == "approved" {
		return errors.New("approved visa applications cannot be deleted")
	}

	if err := vs.Repo.DeleteVisa(id); err != nil {
		return fmt.Errorf("failed to delete visa application: %w", err)
	}
	return nil
}

// GetApprovedVisas retrieves all approved visa applications
func (vs *VisaService) GetApprovedVisas() ([]models.VisaApplication, error) {
	visas, err := vs.Repo.GetApprovedVisas()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve approved visas: %w", err)
	}
	return visas, nil
}

// GetPendingVisas retrieves all pending visa applications
func (vs *VisaService) GetPendingVisas() ([]models.VisaApplication, error) {
	visas, err := vs.Repo.GetPendingVisas()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pending visas: %w", err)
	}
	return visas, nil
}

// GetRejectedVisas retrieves all rejected visa applications
func (vs *VisaService) GetRejectedVisas() ([]models.VisaApplication, error) {
	visas, err := vs.Repo.GetRejectedVisas()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rejected visas: %w", err)
	}
	return visas, nil
}

// ApproveVisa approves a visa application (admin only)
func (vs *VisaService) ApproveVisa(id uint) error {
	if id == 0 {
		return errors.New("invalid visa ID")
	}

	visa, err := vs.Repo.GetVisaById(id)
	if err != nil {
		return fmt.Errorf("visa application not found: %w", err)
	}

	// Check if visa is pending
	if visa.Status != "pending" {
		return fmt.Errorf("only pending visa applications can be approved (current status: %s)", visa.Status)
	}

	visa.Status = "approved"
	if err := vs.Repo.UpdateVisa(visa); err != nil {
		return fmt.Errorf("failed to approve visa application: %w", err)
	}
	return nil
}

// RejectVisa rejects a visa application (admin only)
func (vs *VisaService) RejectVisa(id uint) error {
	if id == 0 {
		return errors.New("invalid visa ID")
	}

	visa, err := vs.Repo.GetVisaById(id)
	if err != nil {
		return fmt.Errorf("visa application not found: %w", err)
	}

	// Check if visa is pending
	if visa.Status != "pending" {
		return fmt.Errorf("only pending visa applications can be rejected (current status: %s)", visa.Status)
	}

	visa.Status = "rejected"
	if err := vs.Repo.UpdateVisa(visa); err != nil {
		return fmt.Errorf("failed to reject visa application: %w", err)
	}
	return nil
}

// GetVisasByDestination retrieves all visa applications for a specific destination
func (vs *VisaService) GetVisasByDestination(destination string) ([]models.VisaApplication, error) {
	if destination == "" {
		return nil, errors.New("destination is required")
	}

	visas, err := vs.Repo.FindVisaByDestination(destination)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve visas for destination %s: %w", destination, err)
	}
	return visas, nil
}

// GetVisasByNationality retrieves all visa applications for a specific nationality
func (vs *VisaService) GetVisasByNationality(nationality string) ([]models.VisaApplication, error) {
	if nationality == "" {
		return nil, errors.New("nationality is required")
	}

	visas, err := vs.Repo.FindVisaByNationality(nationality)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve visas for nationality %s: %w", nationality, err)
	}
	return visas, nil
}

// GetVisasByType retrieves all visa applications for a specific visa type
func (vs *VisaService) GetVisasByType(visaType string) ([]models.VisaApplication, error) {
	if visaType == "" {
		return nil, errors.New("visa type is required")
	}

	visas, err := vs.Repo.FindVisaByVisaType(visaType)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve visas for type %s: %w", visaType, err)
	}
	return visas, nil
}

// GetRecentVisas retrieves the most recent visa applications
func (vs *VisaService) GetRecentVisas(limit int) ([]models.VisaApplication, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	visas, err := vs.Repo.GetRecentVisas(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recent visas: %w", err)
	}
	return visas, nil
}

// GetVisaStatistics returns statistics about visa applications
func (vs *VisaService) GetVisaStatistics() (map[string]int, error) {
	allVisas, err := vs.Repo.GetAllVisas()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve visa statistics: %w", err)
	}

	stats := map[string]int{
		"total":    len(allVisas),
		"pending":  0,
		"approved": 0,
		"rejected": 0,
	}

	for _, visa := range allVisas {
		switch visa.Status {
		case "pending":
			stats["pending"]++
		case "approved":
			stats["approved"]++
		case "rejected":
			stats["rejected"]++
		}
	}

	return stats, nil
}

// BulkApproveVisas approves multiple visa applications (admin only)
func (vs *VisaService) BulkApproveVisas(visaIds []uint) error {
	if len(visaIds) == 0 {
		return errors.New("no visa IDs provided")
	}

	var failedIds []uint
	for _, id := range visaIds {
		if err := vs.ApproveVisa(id); err != nil {
			failedIds = append(failedIds, id)
		}
	}

	if len(failedIds) > 0 {
		return fmt.Errorf("failed to approve visa IDs: %v", failedIds)
	}
	return nil
}

// BulkRejectVisas rejects multiple visa applications (admin only)
func (vs *VisaService) BulkRejectVisas(visaIds []uint) error {
	if len(visaIds) == 0 {
		return errors.New("no visa IDs provided")
	}

	var failedIds []uint
	for _, id := range visaIds {
		if err := vs.RejectVisa(id); err != nil {
			failedIds = append(failedIds, id)
		}
	}

	if len(failedIds) > 0 {
		return fmt.Errorf("failed to reject visa IDs: %v", failedIds)
	}
	return nil
}

// validateVisaApplication validates visa application data
func (vs *VisaService) validateVisaApplication(visa *models.VisaApplication) error {
	// Validate user ID
	if visa.UserID == 0 {
		return errors.New("user ID is required")
	}

	// Validate visa type
	if strings.TrimSpace(visa.VisaType) == "" {
		return errors.New("visa type is required")
	}
	if len(visa.VisaType) < 2 || len(visa.VisaType) > 50 {
		return errors.New("visa type must be between 2 and 50 characters")
	}

	// Validate destination
	if strings.TrimSpace(visa.Destination) == "" {
		return errors.New("destination is required")
	}
	if len(visa.Destination) < 2 || len(visa.Destination) > 100 {
		return errors.New("destination must be between 2 and 100 characters")
	}

	// Validate nationality
	if strings.TrimSpace(visa.Nationality) == "" {
		return errors.New("nationality is required")
	}
	if len(visa.Nationality) < 2 || len(visa.Nationality) > 50 {
		return errors.New("nationality must be between 2 and 50 characters")
	}

	// Validate passport number
	if strings.TrimSpace(visa.PassportNumber) == "" {
		return errors.New("passport number is required")
	}
	if len(visa.PassportNumber) < 6 || len(visa.PassportNumber) > 20 {
		return errors.New("passport number must be between 6 and 20 characters")
	}

	// Validate travel date
	if strings.TrimSpace(visa.TravelDate) == "" {
		return errors.New("travel date is required")
	}

	// Optional: Validate travel date format and that it's in the future
	if _, err := time.Parse("2006-01-02", visa.TravelDate); err != nil {
		return errors.New("invalid travel date format. Use YYYY-MM-DD")
	}

	return nil
}
