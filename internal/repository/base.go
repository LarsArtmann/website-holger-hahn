// Package repository defines data access interfaces and contracts for the portfolio website.
// It provides base repository functionality, common validation logic, and shared utilities
// for implementing repository patterns across different storage backends.
package repository

import (
	"context"
	"fmt"

	"holger-hahn-website/internal/domain"
)

// BaseRepository provides common repository functionality.
type BaseRepository[T Entity] struct {
	entityName string
}

// NewBaseRepository creates a new base repository instance.
func NewBaseRepository[T Entity](entityName string) *BaseRepository[T] {
	return &BaseRepository[T]{
		entityName: entityName,
	}
}

// ValidateEntity performs common entity validation.
func (b *BaseRepository[T]) ValidateEntity(ctx context.Context, entity T) error {
	if entity.GetID() == "" {
		return domain.ErrInvalidInput(b.entityName + " ID cannot be empty")
	}

	return nil
}

// ValidateID validates that an ID is not empty.
func (b *BaseRepository[T]) ValidateID(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput(b.entityName + " ID cannot be empty")
	}

	return nil
}

// NotFoundError creates a standardized not found error.
func (b *BaseRepository[T]) NotFoundError(id string) error {
	return domain.ErrNotFound(fmt.Sprintf("%s with ID %s not found", b.entityName, id))
}

// FilterValidation provides common filter validation.
type FilterValidation struct {
	Limit  *int
	Offset *int
}

// ValidateFilter validates common filter parameters.
func (b *BaseRepository[T]) ValidateFilter(ctx context.Context, filter FilterValidation) error {
	if filter.Limit != nil && *filter.Limit < 0 {
		return domain.ErrInvalidInput("limit cannot be negative")
	}

	if filter.Offset != nil && *filter.Offset < 0 {
		return domain.ErrInvalidInput("offset cannot be negative")
	}

	if filter.Limit != nil && *filter.Limit > 1000 {
		return domain.ErrInvalidInput("limit cannot exceed 1000")
	}

	return nil
}

// GetDefaultLimit returns the default limit for pagination.
func (b *BaseRepository[T]) GetDefaultLimit() int {
	return 50
}

// GetDefaultOffset returns the default offset for pagination.
func (b *BaseRepository[T]) GetDefaultOffset() int {
	return 0
}
