package domain

import (
	"fmt"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    string    `json:"code,omitempty"`
}

// ErrorType represents the type of domain error
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeNotFound   ErrorType = "not_found"
	ErrorTypeConflict   ErrorType = "conflict"
	ErrorTypeInternal   ErrorType = "internal"
)

// Error implements the error interface
func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// ErrInvalidInput creates a validation error
func ErrInvalidInput(message string) error {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// ErrNotFound creates a not found error
func ErrNotFound(resource string) error {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// ErrConflict creates a conflict error
func ErrConflict(message string) error {
	return &DomainError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

// ErrInternal creates an internal error
func ErrInternal(message string) error {
	return &DomainError{
		Type:    ErrorTypeInternal,
		Message: message,
	}
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeValidation
	}
	return false
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeConflict
	}
	return false
}

// IsInternalError checks if the error is an internal error
func IsInternalError(err error) bool {
	if domainErr, ok := err.(*DomainError); ok {
		return domainErr.Type == ErrorTypeInternal
	}
	return false
}