// Package application contains the business logic and use cases for the portfolio website.
// It orchestrates domain entities and coordinates with infrastructure services
// to implement the core business workflows.
package application

import (
	"context"
	"fmt"

	"holger-hahn-website/internal/domain"
)

// ContactService handles contact form business logic.
type ContactService struct {
	contactRepo domain.ContactRepository
	emailSvc    domain.EmailService
	logger      domain.LoggingService
}

// NewContactService creates a new contact service.
func NewContactService(
	contactRepo domain.ContactRepository,
	emailSvc domain.EmailService,
	logger domain.LoggingService,
) *ContactService {
	return &ContactService{
		contactRepo: contactRepo,
		emailSvc:    emailSvc,
		logger:      logger,
	}
}

// SubmitContactForm handles a new contact form submission.
func (s *ContactService) SubmitContactForm(ctx context.Context, req ContactFormRequest) (*ContactFormResponse, error) {
	// Create domain entity with validation.
	contact, err := domain.NewContact(req.Name, req.Company, req.Email, req.Project, "")
	if err != nil {
		s.logger.Error(ctx, "Failed to create contact", err, map[string]interface{}{
			"name":    req.Name,
			"email":   req.Email,
			"company": req.Company,
		})

		return nil, fmt.Errorf("%w: %w", domain.ErrValidationFailed, err)
	}

	// Save to repository.
	if err := s.contactRepo.Save(ctx, contact); err != nil {
		s.logger.Error(ctx, "Failed to save contact", err, map[string]interface{}{
			"contact_id": contact.ID,
			"email":      contact.Email,
		})

		return nil, fmt.Errorf("%w: %w", domain.ErrSaveContact, err)
	}

	s.logger.Info(ctx, "Contact saved successfully", map[string]interface{}{
		"contact_id": contact.ID,
		"email":      contact.Email,
		"company":    contact.Company,
	})

	// Send notification email (async in production).
	if err := s.emailSvc.SendContactNotification(ctx, contact); err != nil {
		s.logger.Error(ctx, "Failed to send notification email", err, map[string]interface{}{
			"contact_id": contact.ID,
			"email":      contact.Email,
		})
		// Don't fail the request if email fails - log and continue.
	}

	// Send confirmation email to contact.
	if err := s.emailSvc.SendConfirmationEmail(ctx, contact); err != nil {
		s.logger.Error(ctx, "Failed to send confirmation email", err, map[string]interface{}{
			"contact_id": contact.ID,
			"email":      contact.Email,
		})
		// Don't fail the request if email fails - log and continue.
	}

	// Mark as processed.
	contact.MarkAsProcessed()

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		s.logger.Error(ctx, "Failed to mark contact as processed", err, map[string]interface{}{
			"contact_id": contact.ID,
		})
		// Don't fail the request - it's already saved.
	}

	return &ContactFormResponse{
		ID:      contact.ID,
		Message: "Thank you for your message! We'll get back to you within 24 hours.",
		Success: true,
	}, nil
}

// GetContact retrieves a contact by ID.
func (s *ContactService) GetContact(ctx context.Context, id string) (*Contact, error) {
	contact, err := s.contactRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrFindContact, err)
	}

	return &Contact{
		ID:          contact.ID,
		Name:        contact.Name,
		Company:     contact.Company,
		Email:       contact.Email,
		Project:     contact.Project,
		Status:      contact.Status,
		SubmittedAt: contact.SubmittedAt,
	}, nil
}

// ListContacts retrieves contacts with pagination.
func (s *ContactService) ListContacts(ctx context.Context, status string, limit, offset int) ([]*Contact, error) {
	contacts, err := s.contactRepo.FindAll(ctx, domain.ContactStatus(status), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domain.ErrListContacts, err)
	}

	result := make([]*Contact, len(contacts))

	for i, contact := range contacts {
		result[i] = &Contact{
			ID:          contact.ID,
			Name:        contact.Name,
			Company:     contact.Company,
			Email:       contact.Email,
			Project:     contact.Project,
			Status:      contact.Status,
			SubmittedAt: contact.SubmittedAt,
		}
	}

	return result, nil
}

// DTOs for application layer.

// ContactFormRequest represents the request payload for contact form submissions,
// containing all required and optional fields with validation rules.
type ContactFormRequest struct {
	Name    string `binding:"required,min=2,max=100" json:"name"`
	Company string `binding:"max=100" json:"company"`
	Email   string `binding:"required,email" json:"email"`
	Project string `binding:"required,min=10,max=2000" json:"project"`
}

// ContactFormResponse represents the response payload after contact form submission,
// including success status and relevant messaging for client feedback.
type ContactFormResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// Contact represents the application layer contact entity for API responses,
// providing a simplified view of contact data for client consumption.
type Contact struct {
	SubmittedAt interface{} `json:"submitted_at"`
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Company     string      `json:"company,omitempty"`
	Email       string      `json:"email"`
	Project     string      `json:"project"`
	Status      string      `json:"status"`
}
