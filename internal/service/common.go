package service

import (
	"context"
	"fmt"
	"strings"

	"holger-hahn-website/internal/domain"
)

// BaseService provides common functionality for all services.
type BaseService struct{}

// NewBaseService creates a new base service.
func NewBaseService() *BaseService {
	return &BaseService{}
}

// ValidateID validates that an ID is not empty.
func (b *BaseService) ValidateID(id, entityType string) error {
	if id == "" {
		return domain.ErrInvalidInput(fmt.Sprintf("%s ID cannot be empty", entityType))
	}
	return nil
}

// ValidateNonEmpty validates that a string field is not empty after trimming.
func (b *BaseService) ValidateNonEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return domain.ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	return nil
}

// NormalizeString trims whitespace from a string.
func (b *BaseService) NormalizeString(value string) string {
	return strings.TrimSpace(value)
}

// WrapRepositoryError wraps repository errors with domain error context.
func (b *BaseService) WrapRepositoryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to %s: %v", operation, err))
}

// HandleNotFound converts repository "not found" errors to domain errors.
func (b *BaseService) HandleNotFound(err error, entityType string) error {
	if err != nil {
		return domain.ErrNotFound(entityType)
	}
	return nil
}

// HandleConflict creates a conflict error with context.
func (b *BaseService) HandleConflict(message string) error {
	return domain.ErrConflict(message)
}

// Validator provides common validation operations.
type Validator struct {
	*BaseService
}

// NewValidator creates a new validator instance.
func NewValidator() *Validator {
	return &Validator{
		BaseService: NewBaseService(),
	}
}

// ValidateEntity performs common entity validation.
type ValidatableEntity interface {
	Validate() error
}

// ValidateAndNormalize validates and normalizes common request fields.
func (v *Validator) ValidateAndNormalize(ctx context.Context, entity ValidatableEntity) error {
	return entity.Validate()
}

// Repository provides common repository operations interface.
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}

// CRUDOperations provides common CRUD operations.
type CRUDOperations[T ValidatableEntity] struct {
	*BaseService
	repo       Repository[T]
	entityType string
}

// NewCRUDOperations creates new CRUD operations helper.
func NewCRUDOperations[T ValidatableEntity](repo Repository[T], entityType string) *CRUDOperations[T] {
	return &CRUDOperations[T]{
		BaseService: NewBaseService(),
		repo:        repo,
		entityType:  entityType,
	}
}

// Get retrieves an entity by ID with common error handling.
func (c *CRUDOperations[T]) Get(ctx context.Context, id string) (*T, error) {
	if err := c.ValidateID(id, c.entityType); err != nil {
		return nil, err
	}

	entity, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, c.HandleNotFound(err, c.entityType)
	}

	return entity, nil
}

// Create creates a new entity with validation and error handling.
func (c *CRUDOperations[T]) Create(ctx context.Context, entity *T) error {
	if err := (*entity).Validate(); err != nil {
		return err
	}

	if err := c.repo.Create(ctx, entity); err != nil {
		return c.WrapRepositoryError(fmt.Sprintf("create %s", c.entityType), err)
	}

	return nil
}

// Update updates an entity with validation and error handling.
func (c *CRUDOperations[T]) Update(ctx context.Context, entity *T) error {
	if err := (*entity).Validate(); err != nil {
		return err
	}

	if err := c.repo.Update(ctx, entity); err != nil {
		return c.WrapRepositoryError(fmt.Sprintf("update %s", c.entityType), err)
	}

	return nil
}

// Delete deletes an entity by ID with error handling.
func (c *CRUDOperations[T]) Delete(ctx context.Context, id string) error {
	if err := c.ValidateID(id, c.entityType); err != nil {
		return err
	}

	// Check if entity exists first
	_, err := c.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := c.repo.Delete(ctx, id); err != nil {
		return c.WrapRepositoryError(fmt.Sprintf("delete %s", c.entityType), err)
	}

	return nil
}

// ServiceConstructor provides a standard constructor pattern.
type ServiceConstructor[T any] struct {
	serviceName string
}

// NewServiceConstructor creates a new service constructor helper.
func NewServiceConstructor[T any](serviceName string) *ServiceConstructor[T] {
	return &ServiceConstructor[T]{
		serviceName: serviceName,
	}
}

// ValidateConstructorArgs validates constructor arguments are not nil.
func (s *ServiceConstructor[T]) ValidateConstructorArgs(args ...interface{}) error {
	for i, arg := range args {
		if arg == nil {
			return domain.ErrInvalidInput(fmt.Sprintf("%s constructor argument %d cannot be nil", s.serviceName, i+1))
		}
	}
	return nil
}
