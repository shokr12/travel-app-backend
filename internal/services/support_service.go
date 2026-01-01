// services/support_service.go
package services

import (
	"Visa/internal/repos"
	"Visa/models"
	"errors"
	"fmt"
	"strings"
)

type SupportService struct {
	Repo *repos.SupportRepo
}

func NewSupportService(supportRepo *repos.SupportRepo) *SupportService {
	return &SupportService{Repo: supportRepo}
}

// GetAllTickets retrieves all support tickets
func (ss *SupportService) GetAllTickets() ([]models.SupportTicket, error) {
	tickets, err := ss.Repo.GetAllTickets()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets: %w", err)
	}
	return tickets, nil
}

// CreateSupportTicket creates a new support ticket
func (ss *SupportService) CreateSupportTicket(ticket *models.SupportTicket) error {
	if ticket == nil {
		return errors.New("ticket data is required")
	}

	// Validate ticket data
	if ticket.UserID == 0 {
		return errors.New("user ID is required")
	}
	if strings.TrimSpace(ticket.Subject) == "" {
		return errors.New("ticket subject is required")
	}
	if len(ticket.Subject) < 5 {
		return errors.New("ticket subject must be at least 5 characters")
	}
	if len(ticket.Subject) > 200 {
		return errors.New("ticket subject must not exceed 200 characters")
	}
	if strings.TrimSpace(ticket.Message) == "" {
		return errors.New("ticket message is required")
	}
	if len(ticket.Message) < 10 {
		return errors.New("ticket message must be at least 10 characters")
	}
	if len(ticket.Message) > 2000 {
		return errors.New("ticket message must not exceed 2000 characters")
	}

	// Set default status if not provided
	if ticket.Status == "" {
		ticket.Status = "open"
	}

	// Validate status
	validStatuses := map[string]bool{
		"open":        true,
		"in_progress": true,
		"resolved":    true,
		"closed":      true,
	}
	if !validStatuses[ticket.Status] {
		return errors.New("invalid ticket status. Must be: open, in_progress, resolved, or closed")
	}

	if err := ss.Repo.CreateTicket(ticket); err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	return nil
}

// GetSupportTicketById retrieves a support ticket by its ID
func (ss *SupportService) GetSupportTicketById(id uint) (*models.SupportTicket, error) {
	if id == 0 {
		return nil, errors.New("invalid ticket ID")
	}

	ticket, err := ss.Repo.GetTicketById(id)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}
	return ticket, nil
}

// GetSupportTicketsByUserId retrieves all support tickets for a specific user
func (ss *SupportService) GetSupportTicketsByUserId(userId uint) ([]models.SupportTicket, error) {
	if userId == 0 {
		return nil, errors.New("invalid user ID")
	}

	tickets, err := ss.Repo.GetTicketsByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user tickets: %w", err)
	}
	return tickets, nil
}

// UpdateSupportTicket updates an existing support ticket
func (ss *SupportService) UpdateSupportTicket(ticket *models.SupportTicket) error {
	if ticket == nil {
		return errors.New("ticket data is required")
	}
	if ticket.ID == 0 {
		return errors.New("ticket ID is required")
	}

	// Verify ticket exists
	existingTicket, err := ss.Repo.GetTicketById(ticket.ID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Validate updated fields if provided
	if ticket.Subject != "" {
		if len(ticket.Subject) < 5 {
			return errors.New("ticket subject must be at least 5 characters")
		}
		if len(ticket.Subject) > 200 {
			return errors.New("ticket subject must not exceed 200 characters")
		}
	}

	if ticket.Message != "" {
		if len(ticket.Message) < 10 {
			return errors.New("ticket message must be at least 10 characters")
		}
		if len(ticket.Message) > 2000 {
			return errors.New("ticket message must not exceed 2000 characters")
		}
	}

	// Validate status if being updated
	if ticket.Status != "" && ticket.Status != existingTicket.Status {
		validStatuses := map[string]bool{
			"open":        true,
			"in_progress": true,
			"resolved":    true,
			"closed":      true,
		}
		if !validStatuses[ticket.Status] {
			return errors.New("invalid ticket status. Must be: open, in_progress, resolved, or closed")
		}

		// Prevent reopening closed tickets
		if existingTicket.Status == "closed" && ticket.Status != "closed" {
			return errors.New("cannot reopen a closed ticket")
		}
	}

	if err := ss.Repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to update ticket: %w", err)
	}
	return nil
}

// DeleteSupportTicket deletes a support ticket
func (ss *SupportService) DeleteSupportTicket(id uint) error {
	if id == 0 {
		return errors.New("invalid ticket ID")
	}

	// Verify ticket exists
	_, err := ss.Repo.GetTicketById(id)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if err := ss.Repo.DeleteTicket(id); err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	return nil
}

// GetTicketsByStatus retrieves all tickets with a specific status
func (ss *SupportService) GetTicketsByStatus(status string) ([]models.SupportTicket, error) {
	if status == "" {
		return nil, errors.New("status is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"open":        true,
		"in_progress": true,
		"resolved":    true,
		"closed":      true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status. Must be: open, in_progress, resolved, or closed")
	}

	tickets, err := ss.Repo.FindByStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets with status %s: %w", status, err)
	}
	return tickets, nil
}

// GetOpenTickets retrieves all open support tickets
func (ss *SupportService) GetOpenTickets() ([]models.SupportTicket, error) {
	return ss.GetTicketsByStatus("open")
}

// GetInProgressTickets retrieves all in-progress support tickets
func (ss *SupportService) GetInProgressTickets() ([]models.SupportTicket, error) {
	return ss.GetTicketsByStatus("in_progress")
}

// GetResolvedTickets retrieves all resolved support tickets
func (ss *SupportService) GetResolvedTickets() ([]models.SupportTicket, error) {
	return ss.GetTicketsByStatus("resolved")
}

// GetClosedTickets retrieves all closed support tickets
func (ss *SupportService) GetClosedTickets() ([]models.SupportTicket, error) {
	return ss.GetTicketsByStatus("closed")
}

// AssignTicket assigns a ticket to in_progress status (admin action)
func (ss *SupportService) AssignTicket(ticketId uint) error {
	if ticketId == 0 {
		return errors.New("invalid ticket ID")
	}

	ticket, err := ss.Repo.GetTicketById(ticketId)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if ticket.Status != "open" {
		return errors.New("only open tickets can be assigned")
	}

	ticket.Status = "in_progress"
	if err := ss.Repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to assign ticket: %w", err)
	}
	return nil
}

// ResolveTicket marks a ticket as resolved (admin action)
func (ss *SupportService) ResolveTicket(ticketId uint) error {
	if ticketId == 0 {
		return errors.New("invalid ticket ID")
	}

	ticket, err := ss.Repo.GetTicketById(ticketId)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if ticket.Status == "closed" {
		return errors.New("cannot resolve a closed ticket")
	}

	ticket.Status = "resolved"
	if err := ss.Repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to resolve ticket: %w", err)
	}
	return nil
}

// CloseTicket closes a support ticket (admin action)
func (ss *SupportService) CloseTicket(ticketId uint) error {
	if ticketId == 0 {
		return errors.New("invalid ticket ID")
	}

	ticket, err := ss.Repo.GetTicketById(ticketId)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	ticket.Status = "closed"
	if err := ss.Repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to close ticket: %w", err)
	}
	return nil
}

// ReopenTicket reopens a resolved ticket (user action)
func (ss *SupportService) ReopenTicket(ticketId uint, userId uint) error {
	if ticketId == 0 {
		return errors.New("invalid ticket ID")
	}
	if userId == 0 {
		return errors.New("invalid user ID")
	}

	ticket, err := ss.Repo.GetTicketById(ticketId)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Verify ticket belongs to user
	if ticket.UserID != userId {
		return errors.New("you can only reopen your own tickets")
	}

	// Only resolved tickets can be reopened
	if ticket.Status != "resolved" {
		return errors.New("only resolved tickets can be reopened")
	}

	ticket.Status = "open"
	if err := ss.Repo.UpdateTicket(ticket); err != nil {
		return fmt.Errorf("failed to reopen ticket: %w", err)
	}
	return nil
}	