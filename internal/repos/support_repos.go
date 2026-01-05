// repos/support_repo.go
package repos

import (
	"Visa/models"

	"gorm.io/gorm"
)

type SupportRepo struct {
	db *gorm.DB
}

func NewSupportRepo(db *gorm.DB) *SupportRepo {
	return &SupportRepo{db: db}
}

// GetAllTickets retrieves all support tickets from the database
func (sr *SupportRepo) GetAllTickets() ([]models.SupportTicket, error) {
	var tickets []models.SupportTicket
	if err := sr.db.Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// GetTicketById retrieves a support ticket by its ID
func (sr *SupportRepo) GetTicketById(id uint) (*models.SupportTicket, error) {
	var ticket models.SupportTicket
	if err := sr.db.First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

// CreateTicket creates a new support ticket in the database
func (sr *SupportRepo) CreateTicket(ticket *models.SupportTicket) error {
	return sr.db.Create(ticket).Error
}

// UpdateTicket updates an existing support ticket
func (sr *SupportRepo) UpdateTicket(ticket *models.SupportTicket) error {
	return sr.db.Save(ticket).Error
}

// DeleteTicket deletes a support ticket by its ID
func (sr *SupportRepo) DeleteTicket(id uint) error {
	return sr.db.Delete(&models.SupportTicket{}, id).Error
}

// FindByStatus retrieves all support tickets with a specific status
// Valid statuses: "open", "in_progress", "resolved", "closed"
func (sr *SupportRepo) FindByStatus(status string) ([]models.SupportTicket, error) {
	var tickets []models.SupportTicket
	if err := sr.db.Preload("User").Where("status = ?", status).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// GetTicketsByUserId retrieves all support tickets for a specific user
func (sr *SupportRepo) GetTicketsByUserId(userId uint) ([]models.SupportTicket, error) {
	var tickets []models.SupportTicket
	if err := sr.db.Preload("User").Where("user_id= ?", userId).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// GetOpenTickets retrieves all open support tickets
func (sr *SupportRepo) GetOpenTickets() ([]models.SupportTicket, error) {
	return sr.FindByStatus("open")
}

// GetInProgressTickets retrieves all in-progress support tickets
func (sr *SupportRepo) GetInProgressTickets() ([]models.SupportTicket, error) {
	return sr.FindByStatus("in_progress")
}

// GetResolvedTickets retrieves all resolved support tickets
func (sr *SupportRepo) GetResolvedTickets() ([]models.SupportTicket, error) {
	return sr.FindByStatus("resolved")
}

// GetClosedTickets retrieves all closed support tickets
func (sr *SupportRepo) GetClosedTickets() ([]models.SupportTicket, error) {
	return sr.FindByStatus("closed")
}
