// Package domain provides core business domain entities and value objects for the portfolio website.
// It defines the contact entity with comprehensive validation logic and business rules
// for handling contact form submissions and contact lifecycle management.
package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Contact represents a contact form submission.
type Contact struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Company     string     `json:"company,omitempty"`
	Message     string     `json:"message"`
	Subject     string     `json:"subject,omitempty"`
	Status      string     `json:"status"`
	Source      string     `json:"source"`
	SubmittedAt time.Time  `json:"submitted_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

// ContactStatus represents the status of a contact submission.
type ContactStatus string

const (
	StatusNew      ContactStatus = "new"
	StatusRead     ContactStatus = "read"
	StatusReplied  ContactStatus = "replied"
	StatusArchived ContactStatus = "archived"
)

// NewContact creates a new contact with validation.
func NewContact(name, company, email, message, subject string) (*Contact, error) {
	if err := validateName(name); err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	if err := validateMessage(message); err != nil {
		return nil, fmt.Errorf("invalid message: %w", err)
	}

	return &Contact{
		ID:          generateID(),
		Name:        strings.TrimSpace(name),
		Company:     strings.TrimSpace(company),
		Email:       strings.ToLower(strings.TrimSpace(email)),
		Message:     strings.TrimSpace(message),
		Subject:     strings.TrimSpace(subject),
		Status:      string(StatusNew),
		Source:      "website",
		SubmittedAt: time.Now().UTC(),
	}, nil
}

// MarkAsRead marks the contact as read.
func (c *Contact) MarkAsRead() {
	c.Status = string(StatusRead)
	now := time.Now().UTC()
	c.ProcessedAt = &now
}

// MarkAsReplied marks the contact as replied.
func (c *Contact) MarkAsReplied() {
	c.Status = string(StatusReplied)
	now := time.Now().UTC()
	c.ProcessedAt = &now
}

// IsValid performs domain validation.
func (c *Contact) IsValid() error {
	if err := validateName(c.Name); err != nil {
		return err
	}

	if err := validateEmail(c.Email); err != nil {
		return err
	}

	if err := validateMessage(c.Message); err != nil {
		return err
	}

	return nil
}

// Validation functions.
func validateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return ErrNameTooShort
	}

	if len(name) > 100 {
		return ErrNameTooLong
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrEmailRequired
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrEmailInvalidFormat
	}

	return nil
}

func validateMessage(message string) error {
	message = strings.TrimSpace(message)
	if len(message) < 10 {
		return ErrMessageTooShort
	}

	if len(message) > 2000 {
		return ErrMessageTooLong
	}

	return nil
}

// Simple ID generation - in production, use UUID or similar.
func generateID() string {
	return fmt.Sprintf("contact_%d", time.Now().UnixNano())
}
