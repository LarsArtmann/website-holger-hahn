// Package domain provides core business domain entities and value objects for the portfolio website.
// It defines domain-specific error types and error handling utilities that encapsulate
// business rule violations and provide meaningful error messages for validation failures.
package domain

import (
	"errors"
	"fmt"
)

// DomainError represents a domain-specific error.
type DomainError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    string    `json:"code,omitempty"`
}

// ErrorType represents the type of domain error.
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeNotFound   ErrorType = "not_found"
	ErrorTypeConflict   ErrorType = "conflict"
	ErrorTypeInternal   ErrorType = "internal"
)

// Error implements the error interface.
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Static error variables to avoid dynamic error creation.
var (
	// Validation errors
	ErrNameRequired           = errors.New("name is required")
	ErrNameTooShort           = errors.New("name must be at least 2 characters long")
	ErrNameTooLong            = errors.New("name must be less than 100 characters")
	ErrEmailRequired          = errors.New("email is required")
	ErrEmailInvalidFormat     = errors.New("invalid email format")
	ErrProjectTooShort        = errors.New("project description must be at least 10 characters long")
	ErrProjectTooLong         = errors.New("project description must be less than 2000 characters")
	ErrInvalidTechnologyLevel = errors.New("invalid technology level")
	ErrTechnologyNameEmpty    = errors.New("technology name cannot be empty")
	ErrTechnologyIDEmpty      = errors.New("technology ID cannot be empty")
	ErrCategoryEmpty          = errors.New("category cannot be empty")

	// Not found errors
	ErrContactNotFound    = errors.New("contact not found")
	ErrTechnologyNotFound = errors.New("technology not found")
	ErrExperienceNotFound = errors.New("experience not found")
	ErrServiceNotFound    = errors.New("service not found")

	// Conflict errors
	ErrTechnologyExists = errors.New("technology already exists")
	ErrContactExists    = errors.New("contact already exists")

	// Repository errors
	ErrContactNil     = errors.New("contact cannot be nil")
	ErrIDEmpty        = errors.New("id cannot be empty")
	ErrInvalidContact = errors.New("invalid contact")

	// Service errors
	ErrCreateTechnology = errors.New("failed to create technology")
	ErrUpdateTechnology = errors.New("failed to update technology")
	ErrDeleteTechnology = errors.New("failed to delete technology")
	ErrListTechnologies = errors.New("failed to list technologies")
	ErrSaveContact      = errors.New("failed to save contact")
	ErrFindContact      = errors.New("failed to find contact")
	ErrListContacts     = errors.New("failed to list contacts")
	ErrValidationFailed = errors.New("validation failed")
	ErrTechByCategory   = errors.New("failed to get technologies by category")
	ErrTechByLevel      = errors.New("failed to get technologies by level")

	// Handler errors
	ErrLoadExperiences  = errors.New("failed to load experiences")
	ErrLoadServices     = errors.New("failed to load services")
	ErrLoadTechnologies = errors.New("failed to load technologies")
	ErrRenderTemplate   = errors.New("failed to render template")
)

// ErrInvalidInput creates a validation error with context.
func ErrInvalidInput(message string) error {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// ErrNotFound creates a not found error with context.
func ErrNotFound(resource string) error {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: resource + " not found",
	}
}

// ErrConflict creates a conflict error with context.
func ErrConflict(message string) error {
	return &DomainError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

// ErrInternal creates an internal error with context.
func ErrInternal(message string) error {
	return &DomainError{
		Type:    ErrorTypeInternal,
		Message: message,
	}
}

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeValidation
	}
	return false
}

// IsNotFoundError checks if the error is a not found error.
func IsNotFoundError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsConflictError checks if the error is a conflict error.
func IsConflictError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeConflict
	}
	return false
}

// IsInternalError checks if the error is an internal error.
func IsInternalError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeInternal
	}
	return false
}
