package service

import (
	"crypto/rand"
	"fmt"
	"time"
)

// like github.com/google/uuid or github.com/matoous/go-nanoid.
func generateID() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)

	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Sprintf("%d-error", timestamp)
	}

	return fmt.Sprintf("%d-%x", timestamp, randomBytes)
}

// ServiceError represents a service layer error with context.
type ServiceError struct {
	Err       error
	Operation string
}

// Error implements the error interface.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("service error in %s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error.
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// WrapError wraps an error with service context.
func WrapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	return &ServiceError{
		Operation: operation,
		Err:       err,
	}
}

// Pagination represents pagination parameters.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// NewPagination creates a new pagination with defaults.
func NewPagination(limit, offset int) *Pagination {
	if limit <= 0 {
		limit = 10 // default limit
	}

	if limit > 100 {
		limit = 100 // max limit
	}

	if offset < 0 {
		offset = 0
	}

	return &Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

// Ordering represents ordering parameters.
type Ordering struct {
	OrderBy  string `json:"order_by"`
	OrderDir string `json:"order_dir"`
}

// NewOrdering creates a new ordering with validation.
func NewOrdering(orderBy, orderDir string) *Ordering {
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "asc" // default to ascending
	}

	return &Ordering{
		OrderBy:  orderBy,
		OrderDir: orderDir,
	}
}

// ValidateOrderBy checks if the order by field is valid for a given entity type.
func ValidateOrderBy(entityType, orderBy string) bool {
	validFields := map[string][]string{
		"technology": {"name", "category", "level", "created_at", "updated_at"},
		"experience": {"start_date", "end_date", "company_name", "position", "created_at", "updated_at"},
		"service":    {"name", "category", "created_at", "updated_at"},
	}

	fields, exists := validFields[entityType]
	if !exists {
		return false
	}

	for _, field := range fields {
		if field == orderBy {
			return true
		}
	}

	return false
}
