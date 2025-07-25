package repository

import (
	"context"
	"time"

	"holger-hahn-website/internal/domain"
)

// ExperienceRepository defines the interface for experience data access.
type ExperienceRepository interface {
	// Create creates a new experience
	Create(ctx context.Context, experience *domain.Experience) error

	// GetByID retrieves an experience by its ID
	GetByID(ctx context.Context, id string) (*domain.Experience, error)

	// List retrieves all experiences with optional filtering and sorting
	List(ctx context.Context, filter ExperienceFilter) ([]*domain.Experience, error)

	// Update updates an existing experience
	Update(ctx context.Context, experience *domain.Experience) error

	// Delete removes an experience by ID
	Delete(ctx context.Context, id string) error

	// GetCurrent retrieves current experiences (where end_date is null)
	GetCurrent(ctx context.Context) ([]*domain.Experience, error)

	// GetByCompany retrieves experiences by company name
	GetByCompany(ctx context.Context, companyName string) ([]*domain.Experience, error)

	// GetByDateRange retrieves experiences within a date range
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Experience, error)

	// GetWithTechnology retrieves experiences that use a specific technology
	GetWithTechnology(ctx context.Context, technologyName string) ([]*domain.Experience, error)
}

// ExperienceFilter represents filtering options for experiences.
type ExperienceFilter struct {
	CompanyName *string
	Position    *string
	IsRemote    *bool
	IsCurrent   *bool
	StartAfter  *time.Time
	EndBefore   *time.Time
	Technology  *string
	Location    *string
	Limit       *int
	Offset      *int
	OrderBy     *string // "start_date", "end_date", "company_name", "position", "created_at"
	OrderDir    *string // "asc", "desc"
}
