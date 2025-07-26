// Package repository defines data access interfaces and contracts for the portfolio website.
// It provides service repository interface with methods for managing professional services,
// filtering by category, pricing, technology, and active status for service offerings.
package repository

import (
	"context"

	"holger-hahn-website/internal/domain"
)

// ServiceRepository defines the interface for service data access.
// It extends the base Filterable interface with service-specific operations.
type ServiceRepository interface {
	Filterable[*domain.Service, ServiceFilter]

	// GetByName retrieves a service by its name
	GetByName(ctx context.Context, name string) (*domain.Service, error)

	// GetActive retrieves all active services
	GetActive(ctx context.Context) ([]*domain.Service, error)

	// GetByCategory retrieves services by category
	GetByCategory(ctx context.Context, category domain.ServiceType) ([]*domain.Service, error)

	// GetWithTechnology retrieves services that use a specific technology
	GetWithTechnology(ctx context.Context, technologyName string) ([]*domain.Service, error)

	// GetByPricingType retrieves services by pricing type
	GetByPricingType(ctx context.Context, pricingType domain.PricingType) ([]*domain.Service, error)
}

// ServiceFilter represents filtering options for services.
type ServiceFilter struct {
	Category    *domain.ServiceType
	IsActive    *bool
	Technology  *string
	PricingType *domain.PricingType
	MinPrice    *float64
	MaxPrice    *float64
	Limit       *int
	Offset      *int
	OrderBy     *string // "name", "category", "created_at", "updated_at"
	OrderDir    *string // "asc", "desc"
}
