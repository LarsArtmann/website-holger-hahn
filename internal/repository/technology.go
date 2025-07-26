package repository

import (
	"context"

	"holger-hahn-website/internal/domain"
)

// TechnologyRepository defines the interface for technology data access.
// It extends the base Filterable interface with technology-specific operations.
type TechnologyRepository interface {
	Filterable[*domain.Technology, TechnologyFilter]

	// GetByName retrieves a technology by its name
	GetByName(ctx context.Context, name string) (*domain.Technology, error)

	// GetByCategory retrieves technologies by category
	GetByCategory(ctx context.Context, category string) ([]*domain.Technology, error)

	// GetByLevel retrieves technologies by proficiency level
	GetByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error)
}

// TechnologyFilter represents filtering options for technologies.
type TechnologyFilter struct {
	Category *string
	Level    *domain.Level
	Limit    *int
	Offset   *int
	OrderBy  *string // "name", "category", "created_at", "updated_at"
	OrderDir *string // "asc", "desc"
}
