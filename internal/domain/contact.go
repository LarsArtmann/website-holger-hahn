package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Contact represents a contact form submission.
type Contact struct {
	SubmittedAt time.Time  `json:"submitted_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Company     string     `json:"company,omitempty"`
	Email       string     `json:"email"`
	Project     string     `json:"project"`
	Status      string     `json:"status"`
}

// ContactStatus represents the status of a contact submission.
type ContactStatus string

const (
	StatusPending   ContactStatus = "pending"
	StatusProcessed ContactStatus = "processed"
	StatusArchived  ContactStatus = "archived"
)

// NewContact creates a new contact with validation.
func NewContact(name, company, email, project string) (*Contact, error) {
	if err := validateName(name); err != nil {
		return nil, fmt.Errorf("invalid name: %w", err)
	}

	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	if err := validateProject(project); err != nil {
		return nil, fmt.Errorf("invalid project: %w", err)
	}

	return &Contact{
		ID:          generateID(),
		Name:        strings.TrimSpace(name),
		Company:     strings.TrimSpace(company),
		Email:       strings.ToLower(strings.TrimSpace(email)),
		Project:     strings.TrimSpace(project),
		Status:      string(StatusPending),
		SubmittedAt: time.Now().UTC(),
	}, nil
}

// MarkAsProcessed marks the contact as processed.
func (c *Contact) MarkAsProcessed() {
	c.Status = string(StatusProcessed)
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

	if err := validateProject(c.Project); err != nil {
		return err
	}

	return nil
}

// Validation functions.
func validateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if len(name) > 100 {
		return errors.New("name must be less than 100 characters")
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func validateProject(project string) error {
	project = strings.TrimSpace(project)
	if len(project) < 10 {
		return errors.New("project description must be at least 10 characters long")
	}

	if len(project) > 2000 {
		return errors.New("project description must be less than 2000 characters")
	}

	return nil
}

// Simple ID generation - in production, use UUID or similar.
func generateID() string {
	return fmt.Sprintf("contact_%d", time.Now().UnixNano())
}
