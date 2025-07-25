package repository

import (
	"context"

	"holger-hahn-website/internal/domain"
)

// ServiceRepository defines the interface for service data access.
type ServiceRepository interface {
	// Create creates a new service
	Create(ctx context.Context, service *domain.Service) error

	// GetByID retrieves a service by its ID
	GetByID(ctx context.Context, id string) (*domain.Service, error)

	// GetByName retrieves a service by its name
	GetByName(ctx context.Context, name string) (*domain.Service, error)

	// List retrieves all services with optional filtering
	List(ctx context.Context, filter ServiceFilter) ([]*domain.Service, error)

	// Update updates an existing service
	Update(ctx context.Context, service *domain.Service) error

	// Delete removes a service by ID
	Delete(ctx context.Context, id string) error

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
