package service

import (
	"context"
	"fmt"

	"holger-hahn-website/internal/domain"
)

// ErrorHandler provides standardized error handling across services.
type ErrorHandler struct {
	serviceName string
}

// NewErrorHandler creates a new error handler for a service.
func NewErrorHandler(serviceName string) *ErrorHandler {
	return &ErrorHandler{
		serviceName: serviceName,
	}
}

// HandleRepositoryError standardizes repository error handling.
func (e *ErrorHandler) HandleRepositoryError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to %s: %v", operation, err))
}

// HandleCreateError standardizes create operation error handling.
func (e *ErrorHandler) HandleCreateError(entityType string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to create %s: %v", entityType, err))
}

// HandleUpdateError standardizes update operation error handling.
func (e *ErrorHandler) HandleUpdateError(entityType string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to update %s: %v", entityType, err))
}

// HandleDeleteError standardizes delete operation error handling.
func (e *ErrorHandler) HandleDeleteError(entityType string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to delete %s: %v", entityType, err))
}

// HandleListError standardizes list operation error handling.
func (e *ErrorHandler) HandleListError(entityType string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to list %s: %v", entityType, err))
}

// HandleGetError standardizes get operation error handling.
func (e *ErrorHandler) HandleGetError(entityType string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrNotFound(entityType)
}

// HandleCustomOperation standardizes custom operation error handling.
func (e *ErrorHandler) HandleCustomOperation(operationName string, err error) error {
	if err == nil {
		return nil
	}
	return domain.ErrInternal(fmt.Sprintf("failed to %s: %v", operationName, err))
}

// ServiceErrorContext provides contextual error information.
type ServiceErrorContext struct {
	ServiceName string
	Operation   string
	EntityType  string
	EntityID    string
	Error       error
}

// ContextualErrorHandler provides error handling with rich context.
type ContextualErrorHandler struct {
	*ErrorHandler
}

// NewContextualErrorHandler creates a new contextual error handler.
func NewContextualErrorHandler(serviceName string) *ContextualErrorHandler {
	return &ContextualErrorHandler{
		ErrorHandler: NewErrorHandler(serviceName),
	}
}

// HandleWithContext handles errors with full context information.
func (c *ContextualErrorHandler) HandleWithContext(ctx context.Context, errorCtx ServiceErrorContext) error {
	if errorCtx.Error == nil {
		return nil
	}

	// Add context to error message
	contextMessage := fmt.Sprintf("[%s.%s]", errorCtx.ServiceName, errorCtx.Operation)
	if errorCtx.EntityType != "" {
		contextMessage += fmt.Sprintf(" %s", errorCtx.EntityType)
	}
	if errorCtx.EntityID != "" {
		contextMessage += fmt.Sprintf(" (ID: %s)", errorCtx.EntityID)
	}

	return domain.ErrInternal(fmt.Sprintf("%s: %v", contextMessage, errorCtx.Error))
}

// RepositoryErrorHandler specializes in handling repository-specific errors.
type RepositoryErrorHandler struct {
	*ContextualErrorHandler
}

// NewRepositoryErrorHandler creates a repository-specific error handler.
func NewRepositoryErrorHandler(serviceName string) *RepositoryErrorHandler {
	return &RepositoryErrorHandler{
		ContextualErrorHandler: NewContextualErrorHandler(serviceName),
	}
}

// HandleNotFound handles repository "not found" errors consistently.
func (r *RepositoryErrorHandler) HandleNotFound(entityType string, id string) error {
	return domain.ErrNotFound(entityType)
}

// HandleConflict handles repository conflict errors consistently.
func (r *RepositoryErrorHandler) HandleConflict(entityType, entityName string) error {
	return domain.ErrConflict(fmt.Sprintf("%s '%s' already exists", entityType, entityName))
}

// HandleAlreadyExists handles "already exists" errors consistently.
func (r *RepositoryErrorHandler) HandleAlreadyExists(entityType, identifier string) error {
	return domain.ErrConflict(fmt.Sprintf("%s '%s' already exists", entityType, identifier))
}

// BusinessLogicErrorHandler handles business logic specific errors.
type BusinessLogicErrorHandler struct {
	*ContextualErrorHandler
}

// NewBusinessLogicErrorHandler creates a business logic error handler.
func NewBusinessLogicErrorHandler(serviceName string) *BusinessLogicErrorHandler {
	return &BusinessLogicErrorHandler{
		ContextualErrorHandler: NewContextualErrorHandler(serviceName),
	}
}

// HandleRelationshipError handles relationship-related errors.
func (b *BusinessLogicErrorHandler) HandleRelationshipError(fromEntity, toEntity, relationshipType string, err error) error {
	if err == nil {
		return nil
	}
	
	return domain.ErrInternal(fmt.Sprintf("failed to establish %s relationship between %s and %s: %v", 
		relationshipType, fromEntity, toEntity, err))
}

// HandleDuplicateRelationship handles duplicate relationship errors.
func (b *BusinessLogicErrorHandler) HandleDuplicateRelationship(relationshipType string) error {
	return domain.ErrConflict(fmt.Sprintf("%s already exists", relationshipType))
}

// HandleInvalidStateTransition handles invalid state transition errors.
func (b *BusinessLogicErrorHandler) HandleInvalidStateTransition(entityType, currentState, targetState string) error {
	return domain.ErrInvalidInput(fmt.Sprintf("cannot transition %s from %s to %s", entityType, currentState, targetState))
}

// ErrorHandlerFactory creates appropriate error handlers for different contexts.
type ErrorHandlerFactory struct{}

// NewErrorHandlerFactory creates a new error handler factory.
func NewErrorHandlerFactory() *ErrorHandlerFactory {
	return &ErrorHandlerFactory{}
}

// CreateForService creates a basic error handler for a service.
func (f *ErrorHandlerFactory) CreateForService(serviceName string) *ErrorHandler {
	return NewErrorHandler(serviceName)
}

// CreateContextualForService creates a contextual error handler for a service.
func (f *ErrorHandlerFactory) CreateContextualForService(serviceName string) *ContextualErrorHandler {
	return NewContextualErrorHandler(serviceName)
}

// CreateRepositoryErrorHandler creates a repository error handler for a service.
func (f *ErrorHandlerFactory) CreateRepositoryErrorHandler(serviceName string) *RepositoryErrorHandler {
	return NewRepositoryErrorHandler(serviceName)
}

// CreateBusinessLogicErrorHandler creates a business logic error handler for a service.
func (f *ErrorHandlerFactory) CreateBusinessLogicErrorHandler(serviceName string) *BusinessLogicErrorHandler {
	return NewBusinessLogicErrorHandler(serviceName)
}

// StandardServiceErrorHandlers provides all common error handlers for a service.
type StandardServiceErrorHandlers struct {
	Basic        *ErrorHandler
	Contextual   *ContextualErrorHandler
	Repository   *RepositoryErrorHandler
	BusinessLogic *BusinessLogicErrorHandler
}

// NewStandardServiceErrorHandlers creates a complete set of error handlers for a service.
func NewStandardServiceErrorHandlers(serviceName string) *StandardServiceErrorHandlers {
	return &StandardServiceErrorHandlers{
		Basic:        NewErrorHandler(serviceName),
		Contextual:   NewContextualErrorHandler(serviceName),
		Repository:   NewRepositoryErrorHandler(serviceName),
		BusinessLogic: NewBusinessLogicErrorHandler(serviceName),
	}
}