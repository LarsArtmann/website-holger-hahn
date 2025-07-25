package domain

import "context"

// ContactRepository defines the interface for contact persistence
type ContactRepository interface {
	// Save stores a contact submission
	Save(ctx context.Context, contact *Contact) error
	
	// FindByID retrieves a contact by ID
	FindByID(ctx context.Context, id string) (*Contact, error)
	
	// FindAll retrieves all contacts with optional filtering
	FindAll(ctx context.Context, status ContactStatus, limit, offset int) ([]*Contact, error)
	
	// Update updates an existing contact
	Update(ctx context.Context, contact *Contact) error
	
	// Delete removes a contact (for GDPR compliance)
	Delete(ctx context.Context, id string) error
	
	// Count returns the total number of contacts
	Count(ctx context.Context, status ContactStatus) (int, error)
}

// EmailService defines the interface for sending emails
type EmailService interface {
	// SendContactNotification sends a notification email about a new contact
	SendContactNotification(ctx context.Context, contact *Contact) error
	
	// SendConfirmationEmail sends a confirmation email to the contact
	SendConfirmationEmail(ctx context.Context, contact *Contact) error
}

// LoggingService defines the interface for structured logging
type LoggingService interface {
	// Info logs an info message with structured data
	Info(ctx context.Context, message string, fields map[string]interface{})
	
	// Error logs an error message with structured data
	Error(ctx context.Context, message string, err error, fields map[string]interface{})
	
	// Warn logs a warning message with structured data
	Warn(ctx context.Context, message string, fields map[string]interface{})
}