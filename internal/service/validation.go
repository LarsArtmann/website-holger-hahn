package service

import (
	"fmt"
	"strings"
	"time"

	"holger-hahn-website/internal/domain"
)

// InputValidator provides common input validation operations.
type InputValidator struct {
	*BaseService
}

// NewInputValidator creates a new input validator.
func NewInputValidator() *InputValidator {
	return &InputValidator{
		BaseService: NewBaseService(),
	}
}

// ValidateStringField validates and normalizes a string field.
type StringFieldValidator struct {
	fieldName string
	value     string
	required  bool
	minLength int
	maxLength int
}

// NewStringField creates a new string field validator.
func NewStringField(fieldName, value string) *StringFieldValidator {
	return &StringFieldValidator{
		fieldName: fieldName,
		value:     value,
	}
}

// Required marks the field as required.
func (s *StringFieldValidator) Required() *StringFieldValidator {
	s.required = true
	return s
}

// MinLength sets minimum length validation.
func (s *StringFieldValidator) MinLength(min int) *StringFieldValidator {
	s.minLength = min
	return s
}

// MaxLength sets maximum length validation.
func (s *StringFieldValidator) MaxLength(max int) *StringFieldValidator {
	s.maxLength = max
	return s
}

// Validate performs the validation and returns normalized value and error.
func (s *StringFieldValidator) Validate() (string, error) {
	// Normalize the value
	normalized := strings.TrimSpace(s.value)

	// Check if required
	if s.required && normalized == "" {
		return "", domain.ErrInvalidInput(fmt.Sprintf("%s cannot be empty", s.fieldName))
	}

	// Check minimum length
	if s.minLength > 0 && len(normalized) > 0 && len(normalized) < s.minLength {
		return "", domain.ErrInvalidInput(fmt.Sprintf("%s must be at least %d characters long", s.fieldName, s.minLength))
	}

	// Check maximum length
	if s.maxLength > 0 && len(normalized) > s.maxLength {
		return "", domain.ErrInvalidInput(fmt.Sprintf("%s must be less than %d characters", s.fieldName, s.maxLength))
	}

	return normalized, nil
}

// CommonRequestValidator validates common request patterns across services.
type CommonRequestValidator struct {
	*InputValidator
}

// NewCommonRequestValidator creates a new common request validator.
func NewCommonRequestValidator() *CommonRequestValidator {
	return &CommonRequestValidator{
		InputValidator: NewInputValidator(),
	}
}

// ValidateCreateRequest validates common create request fields.
func (v *CommonRequestValidator) ValidateCreateRequest(name, description string) (string, string, error) {
	// Validate name
	normalizedName, err := NewStringField("name", name).Required().MinLength(2).MaxLength(100).Validate()
	if err != nil {
		return "", "", err
	}

	// Validate description
	normalizedDescription, err := NewStringField("description", description).MaxLength(2000).Validate()
	if err != nil {
		return "", "", err
	}

	return normalizedName, normalizedDescription, nil
}

// ValidateUpdateRequest validates common update request fields.
func (v *CommonRequestValidator) ValidateUpdateRequest(id string, entityType string) error {
	return v.ValidateID(id, entityType)
}

// ValidateListRequest validates common list/filter request parameters.
func (v *CommonRequestValidator) ValidateListRequest(limit, offset *int, orderBy, orderDir *string, entityType string) error {
	// Validate pagination
	if limit != nil && *limit < 0 {
		return domain.ErrInvalidInput("limit cannot be negative")
	}

	if limit != nil && *limit > 100 {
		return domain.ErrInvalidInput("limit cannot exceed 100")
	}

	if offset != nil && *offset < 0 {
		return domain.ErrInvalidInput("offset cannot be negative")
	}

	// Validate ordering
	if orderDir != nil && *orderDir != "asc" && *orderDir != "desc" {
		return domain.ErrInvalidInput("order direction must be 'asc' or 'desc'")
	}

	if orderBy != nil && !ValidateOrderBy(entityType, *orderBy) {
		return domain.ErrInvalidInput(fmt.Sprintf("invalid order by field '%s' for %s", *orderBy, entityType))
	}

	return nil
}

// EntityExistenceValidator validates entity existence patterns.
type EntityExistenceValidator struct {
	*BaseService
}

// NewEntityExistenceValidator creates a new entity existence validator.
func NewEntityExistenceValidator() *EntityExistenceValidator {
	return &EntityExistenceValidator{
		BaseService: NewBaseService(),
	}
}

// ValidateEntityExists validates that an entity exists before operations.
func (v *EntityExistenceValidator) ValidateEntityExists(entity interface{}, entityType string) error {
	if entity == nil {
		return domain.ErrNotFound(entityType)
	}
	return nil
}

// ValidateEntityNotExists validates that an entity doesn't exist (for conflict detection).
func (v *EntityExistenceValidator) ValidateEntityNotExists(entity interface{}, entityName, entityType string) error {
	if entity != nil {
		return domain.ErrConflict(fmt.Sprintf("%s '%s' already exists", entityType, entityName))
	}
	return nil
}

// RelationshipValidator validates entity relationships.
type RelationshipValidator struct {
	*BaseService
}

// NewRelationshipValidator creates a new relationship validator.
func NewRelationshipValidator() *RelationshipValidator {
	return &RelationshipValidator{
		BaseService: NewBaseService(),
	}
}

// ValidateRelationshipNotExists validates that a relationship doesn't already exist.
func (v *RelationshipValidator) ValidateRelationshipNotExists(exists bool, relationshipType string) error {
	if exists {
		return domain.ErrConflict(fmt.Sprintf("%s already exists", relationshipType))
	}
	return nil
}

// DateValidator validates date-related fields.
type DateValidator struct {
	*BaseService
}

// NewDateValidator creates a new date validator.
func NewDateValidator() *DateValidator {
	return &DateValidator{
		BaseService: NewBaseService(),
	}
}

// ValidateStartEndDates validates that start date is before end date.
func (v *DateValidator) ValidateStartEndDates(startDate time.Time, endDate *time.Time) error {
	if endDate != nil && endDate.Before(startDate) {
		return domain.ErrInvalidInput("end date cannot be before start date")
	}
	return nil
}

// ValidateDateNotInFuture validates that a date is not in the future.
func (v *DateValidator) ValidateDateNotInFuture(date time.Time, fieldName string) error {
	if date.After(time.Now()) {
		return domain.ErrInvalidInput(fmt.Sprintf("%s cannot be in the future", fieldName))
	}
	return nil
}

// CompoundValidator combines multiple validators for complex validation scenarios.
type CompoundValidator struct {
	*CommonRequestValidator
	*EntityExistenceValidator
	*RelationshipValidator
	*DateValidator
}

// NewCompoundValidator creates a new compound validator with all validation capabilities.
func NewCompoundValidator() *CompoundValidator {
	return &CompoundValidator{
		CommonRequestValidator:   NewCommonRequestValidator(),
		EntityExistenceValidator: NewEntityExistenceValidator(),
		RelationshipValidator:    NewRelationshipValidator(),
		DateValidator:            NewDateValidator(),
	}
}
