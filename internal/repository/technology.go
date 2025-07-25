package repository

import (
	"context"
	"holger-hahn-website/internal/domain"
)

// TechnologyRepository defines the interface for technology data access
type TechnologyRepository interface {
	// Create creates a new technology
	Create(ctx context.Context, technology *domain.Technology) error
	
	// GetByID retrieves a technology by its ID
	GetByID(ctx context.Context, id string) (*domain.Technology, error)
	
	// GetByName retrieves a technology by its name
	GetByName(ctx context.Context, name string) (*domain.Technology, error)
	
	// List retrieves all technologies with optional filtering
	List(ctx context.Context, filter TechnologyFilter) ([]*domain.Technology, error)
	
	// Update updates an existing technology
	Update(ctx context.Context, technology *domain.Technology) error
	
	// Delete removes a technology by ID
	Delete(ctx context.Context, id string) error
	
	// GetByCategory retrieves technologies by category
	GetByCategory(ctx context.Context, category string) ([]*domain.Technology, error)
	
	// GetByLevel retrieves technologies by proficiency level
	GetByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error)
}

// TechnologyFilter represents filtering options for technologies
type TechnologyFilter struct {
	Category *string
	Level    *domain.Level
	Limit    *int
	Offset   *int
	OrderBy  *string // "name", "category", "created_at", "updated_at"
	OrderDir *string // "asc", "desc"
}