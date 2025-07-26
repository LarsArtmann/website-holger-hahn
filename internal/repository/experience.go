// Package repository defines data access interfaces and contracts for the portfolio website.
// It provides experience repository interface with methods for managing professional
// experience data, filtering by company, date ranges, and technology associations.
package repository

import (
	"context"
	"time"

	"holger-hahn-website/internal/domain"
)

// ExperienceRepository defines the interface for experience data access.
// It extends the base Filterable interface with experience-specific operations.
type ExperienceRepository interface {
	Filterable[*domain.Experience, ExperienceFilter]

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
