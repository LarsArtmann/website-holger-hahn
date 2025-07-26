// Package service provides business logic and application services for the portfolio website.
// It contains domain services that coordinate business operations, enforce business rules,
// and orchestrate interactions between domain entities and repositories.
package service

import (
	"context"
	"fmt"
	"strings"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// PortfolioService handles business logic for service offerings.
type PortfolioService struct {
	repo     repository.ServiceRepository
	techRepo repository.TechnologyRepository
}

// NewPortfolioService creates a new portfolio service.
func NewPortfolioService(repo repository.ServiceRepository, techRepo repository.TechnologyRepository) *PortfolioService {
	return &PortfolioService{
		repo:     repo,
		techRepo: techRepo,
	}
}

// CreateService creates a new service offering with validation.
func (s *PortfolioService) CreateService(ctx context.Context, req CreateServiceRequest) (*domain.Service, error) {
	// Normalize input
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return nil, domain.ErrInvalidInput("service name cannot be empty")
	}

	if req.Description == "" {
		return nil, domain.ErrInvalidInput("service description cannot be empty")
	}

	// Check if service already exists
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, domain.ErrConflict(fmt.Sprintf("service '%s' already exists", req.Name))
	}

	// Create new service
	service := domain.NewService(req.Name, req.Description, req.Category)
	service.ID = generateID()
	service.Duration = req.Duration

	// Set pricing if provided
	if req.Pricing != nil {
		if err := service.SetPricing(*req.Pricing); err != nil {
			return nil, err
		}
	}

	if err := service.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, service); err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to create service: %v", err))
	}

	return service, nil
}

// GetService retrieves a service by ID.
func (s *PortfolioService) GetService(ctx context.Context, id string) (*domain.Service, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput("service ID cannot be empty")
	}

	service, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound("service")
	}

	return service, nil
}

// ListServices retrieves services with filtering and pagination.
func (s *PortfolioService) ListServices(ctx context.Context, filter ServiceFilter) ([]*domain.Service, error) {
	repoFilter := repository.ServiceFilter{
		Category:    filter.Category,
		IsActive:    filter.IsActive,
		Technology:  filter.Technology,
		PricingType: filter.PricingType,
		MinPrice:    filter.MinPrice,
		MaxPrice:    filter.MaxPrice,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
		OrderBy:     filter.OrderBy,
		OrderDir:    filter.OrderDir,
	}

	services, err := s.repo.List(ctx, repoFilter)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to list services: %v", err))
	}

	return services, nil
}

// GetActiveServices retrieves all active services.
func (s *PortfolioService) GetActiveServices(ctx context.Context) ([]*domain.Service, error) {
	services, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to get active services: %v", err))
	}

	return services, nil
}

// GetServicesByCategory retrieves services by category.
func (s *PortfolioService) GetServicesByCategory(ctx context.Context, category domain.ServiceType) ([]*domain.Service, error) {
	if !category.IsValid() {
		return nil, domain.ErrInvalidInput("invalid service category")
	}

	services, err := s.repo.GetByCategory(ctx, category)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to get services by category: %v", err))
	}

	return services, nil
}

// AddTechnologyToService adds a technology to a service.
func (s *PortfolioService) AddTechnologyToService(ctx context.Context, serviceID, technologyID string) error {
	// Get service
	service, err := s.GetService(ctx, serviceID)
	if err != nil {
		return err
	}

	// Get technology
	technology, err := s.techRepo.GetByID(ctx, technologyID)
	if err != nil {
		return domain.ErrNotFound("technology")
	}

	// Check if technology is already added
	for _, tech := range service.Technologies {
		if tech.ID == technologyID {
			return domain.ErrConflict("technology already added to service")
		}
	}

	// Add technology
	service.AddTechnology(*technology)

	if err := s.repo.Update(ctx, service); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to add technology to service: %v", err))
	}

	return nil
}

// AddDeliverableToService adds a deliverable to a service.
func (s *PortfolioService) AddDeliverableToService(ctx context.Context, serviceID string, deliverable DeliverableRequest) error {
	// Get service
	service, err := s.GetService(ctx, serviceID)
	if err != nil {
		return err
	}

	// Create deliverable
	del := domain.NewDeliverable(deliverable.Name, deliverable.Description, deliverable.Timeline)
	del.ID = generateID()

	// Add deliverable
	service.AddDeliverable(*del)

	if err := s.repo.Update(ctx, service); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to add deliverable to service: %v", err))
	}

	return nil
}

// UpdateServicePricing updates the pricing for a service.
func (s *PortfolioService) UpdateServicePricing(ctx context.Context, serviceID string, pricing domain.PricingInfo) error {
	service, err := s.GetService(ctx, serviceID)
	if err != nil {
		return err
	}

	if err := service.SetPricing(pricing); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, service); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to update service pricing: %v", err))
	}

	return nil
}

// ActivateService activates a service.
func (s *PortfolioService) ActivateService(ctx context.Context, id string) error {
	service, err := s.GetService(ctx, id)
	if err != nil {
		return err
	}

	service.Activate()

	if err := s.repo.Update(ctx, service); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to activate service: %v", err))
	}

	return nil
}

// DeactivateService deactivates a service.
func (s *PortfolioService) DeactivateService(ctx context.Context, id string) error {
	service, err := s.GetService(ctx, id)
	if err != nil {
		return err
	}

	service.Deactivate()

	if err := s.repo.Update(ctx, service); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to deactivate service: %v", err))
	}

	return nil
}

// CreateServiceRequest represents a request to create a service.
type CreateServiceRequest struct {
	Pricing     *domain.PricingInfo `json:"pricing,omitempty"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Category    domain.ServiceType  `json:"category"`
	Duration    string              `json:"duration,omitempty"`
}

// ServiceFilter represents filtering options for the service layer.
type ServiceFilter struct {
	Category    *domain.ServiceType
	IsActive    *bool
	Technology  *string
	PricingType *domain.PricingType
	MinPrice    *float64
	MaxPrice    *float64
	Limit       *int
	Offset      *int
	OrderBy     *string
	OrderDir    *string
}

// DeliverableRequest represents a request to add a deliverable.
type DeliverableRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Timeline    string `json:"timeline,omitempty"`
}
