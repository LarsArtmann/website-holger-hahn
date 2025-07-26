// Package container provides in-memory repository implementations for development and testing.
// These implementations store data in memory and provide the same interface as production
// repositories for local development and unit testing purposes.
package container

import (
	"context"
	"sync"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// InMemoryTechnologyRepository is a simple in-memory implementation for development/testing.
type InMemoryTechnologyRepository struct {
	techs map[string]*domain.Technology
	mu    sync.RWMutex
}

// NewInMemoryTechnologyRepository creates a new in-memory technology repository.
func NewInMemoryTechnologyRepository() repository.TechnologyRepository {
	return &InMemoryTechnologyRepository{
		techs: make(map[string]*domain.Technology),
	}
}

// Create stores a new technology in the in-memory repository.
func (r *InMemoryTechnologyRepository) Create(ctx context.Context, technology *domain.Technology) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.techs[technology.ID] = technology

	return nil
}

// GetByID retrieves a technology by its ID.
func (r *InMemoryTechnologyRepository) GetByID(ctx context.Context, id string) (*domain.Technology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.techs[id]
	if !exists {
		return nil, domain.ErrNotFound("technology")
	}

	return item, nil
}

// GetByName finds a technology by its name.
func (r *InMemoryTechnologyRepository) GetByName(ctx context.Context, name string) (*domain.Technology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tech := range r.techs {
		if tech.Name == name {
			return tech, nil
		}
	}

	return nil, domain.ErrNotFound("technology")
}

// List retrieves technologies matching the provided filter criteria.
func (r *InMemoryTechnologyRepository) List(ctx context.Context, filter repository.TechnologyFilter) ([]*domain.Technology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Technology

	for _, tech := range r.techs {
		if filter.Category != nil && tech.Category != *filter.Category {
			continue
		}

		if filter.Level != nil && tech.Level != *filter.Level {
			continue
		}

		result = append(result, tech)
	}

	return result, nil
}

// Update modifies an existing technology in the in-memory repository.
func (r *InMemoryTechnologyRepository) Update(ctx context.Context, technology *domain.Technology) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.techs[technology.ID] = technology

	return nil
}

// Delete removes a technology from the in-memory repository by its ID.
func (r *InMemoryTechnologyRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.techs, id)

	return nil
}

// GetByCategory retrieves all technologies belonging to the specified category.
func (r *InMemoryTechnologyRepository) GetByCategory(ctx context.Context, category string) ([]*domain.Technology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Technology

	for _, tech := range r.techs {
		if tech.Category == category {
			result = append(result, tech)
		}
	}

	return result, nil
}

// GetByLevel retrieves all technologies at the specified proficiency level.
func (r *InMemoryTechnologyRepository) GetByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Technology

	for _, tech := range r.techs {
		if tech.Level == level {
			result = append(result, tech)
		}
	}

	return result, nil
}

// InMemoryExperienceRepository is a simple in-memory implementation for development/testing.
type InMemoryExperienceRepository struct {
	exps map[string]*domain.Experience
	mu   sync.RWMutex
}

// NewInMemoryExperienceRepository creates a new in-memory experience repository.
func NewInMemoryExperienceRepository() repository.ExperienceRepository {
	return &InMemoryExperienceRepository{
		exps: make(map[string]*domain.Experience),
	}
}

// Create stores a new experience in the in-memory repository.
func (r *InMemoryExperienceRepository) Create(ctx context.Context, experience *domain.Experience) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exps[experience.ID] = experience

	return nil
}

// GetByID retrieves an experience by its ID.
func (r *InMemoryExperienceRepository) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.exps[id]
	if !exists {
		return nil, domain.ErrNotFound("experience")
	}

	return item, nil
}

// List retrieves experiences matching the provided filter criteria.
func (r *InMemoryExperienceRepository) List(ctx context.Context, filter repository.ExperienceFilter) ([]*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Experience

	for _, exp := range r.exps {
		result = append(result, exp)
	}

	return result, nil
}

// Update modifies an existing experience in the in-memory repository.
func (r *InMemoryExperienceRepository) Update(ctx context.Context, experience *domain.Experience) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exps[experience.ID] = experience

	return nil
}

// Delete removes an experience from the in-memory repository by its ID.
func (r *InMemoryExperienceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.exps, id)

	return nil
}

// GetCurrent retrieves all experiences that are currently active.
func (r *InMemoryExperienceRepository) GetCurrent(ctx context.Context) ([]*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Experience

	for _, exp := range r.exps {
		if exp.IsCurrent() {
			result = append(result, exp)
		}
	}

	return result, nil
}

// GetByCompany retrieves all experiences associated with the specified company.
func (r *InMemoryExperienceRepository) GetByCompany(ctx context.Context, companyName string) ([]*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Experience

	for _, exp := range r.exps {
		if exp.CompanyName == companyName {
			result = append(result, exp)
		}
	}

	return result, nil
}

// GetByDateRange retrieves experiences that fall within the specified date range.
func (r *InMemoryExperienceRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Experience

	for _, exp := range r.exps {
		if !exp.StartDate.Before(startDate) && (exp.EndDate == nil || !exp.EndDate.After(endDate)) {
			result = append(result, exp)
		}
	}

	return result, nil
}

// GetWithTechnology retrieves experiences that utilize the specified technology.
func (r *InMemoryExperienceRepository) GetWithTechnology(ctx context.Context, technologyName string) ([]*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Experience

	for _, exp := range r.exps {
		for _, tech := range exp.Technologies {
			if tech.Name == technologyName {
				result = append(result, exp)
				break
			}
		}
	}

	return result, nil
}

// InMemoryServiceRepository is a simple in-memory implementation for development/testing.
type InMemoryServiceRepository struct {
	services map[string]*domain.Service
	mu       sync.RWMutex
}

// NewInMemoryServiceRepository creates a new in-memory service repository.
func NewInMemoryServiceRepository() repository.ServiceRepository {
	return &InMemoryServiceRepository{
		services: make(map[string]*domain.Service),
	}
}

// Create stores a new service in the in-memory repository.
func (r *InMemoryServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[service.ID] = service

	return nil
}

// GetByID retrieves a service by its ID.
func (r *InMemoryServiceRepository) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.services[id]
	if !exists {
		return nil, domain.ErrNotFound("service")
	}

	return item, nil
}

// GetByName finds a service by its name.
func (r *InMemoryServiceRepository) GetByName(ctx context.Context, name string) (*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, svc := range r.services {
		if svc.Name == name {
			return svc, nil
		}
	}

	return nil, domain.ErrNotFound("service")
}

// List retrieves services matching the provided filter criteria.
func (r *InMemoryServiceRepository) List(ctx context.Context, filter repository.ServiceFilter) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Service

	for _, svc := range r.services {
		if filter.Category != nil && svc.Category != *filter.Category {
			continue
		}

		if filter.IsActive != nil && svc.IsActive != *filter.IsActive {
			continue
		}

		result = append(result, svc)
	}

	return result, nil
}

// Update modifies an existing service in the in-memory repository.
func (r *InMemoryServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[service.ID] = service

	return nil
}

// Delete removes a service from the in-memory repository by its ID.
func (r *InMemoryServiceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, id)

	return nil
}

// GetActive retrieves all services that are currently active.
func (r *InMemoryServiceRepository) GetActive(ctx context.Context) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Service

	for _, svc := range r.services {
		if svc.IsActive {
			result = append(result, svc)
		}
	}

	return result, nil
}

// GetByCategory retrieves all services belonging to the specified category.
func (r *InMemoryServiceRepository) GetByCategory(ctx context.Context, category domain.ServiceType) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Service

	for _, svc := range r.services {
		if svc.Category == category {
			result = append(result, svc)
		}
	}

	return result, nil
}

// GetWithTechnology retrieves services that utilize the specified technology.
func (r *InMemoryServiceRepository) GetWithTechnology(ctx context.Context, technologyName string) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Service

	for _, svc := range r.services {
		for _, tech := range svc.Technologies {
			if tech.Name == technologyName {
				result = append(result, svc)
				break
			}
		}
	}

	return result, nil
}

// GetByPricingType retrieves services that use the specified pricing model.
func (r *InMemoryServiceRepository) GetByPricingType(ctx context.Context, pricingType domain.PricingType) ([]*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Service

	for _, svc := range r.services {
		if svc.Pricing != nil && svc.Pricing.Type == pricingType {
			result = append(result, svc)
		}
	}

	return result, nil
}

// InMemoryUnitOfWork is a simple in-memory implementation for development/testing.
type InMemoryUnitOfWork struct{}

// NewInMemoryUnitOfWork creates a new in-memory unit of work.
func NewInMemoryUnitOfWork() repository.UnitOfWork {
	return &InMemoryUnitOfWork{}
}

// Begin starts a new transaction for the unit of work.
func (u *InMemoryUnitOfWork) Begin(ctx context.Context) (repository.Transaction, error) {
	return &InMemoryTransaction{}, nil
}

// InMemoryTransaction is a simple in-memory transaction implementation.
type InMemoryTransaction struct{}

// Technology returns the technology repository for this transaction.
func (t *InMemoryTransaction) Technology() repository.TechnologyRepository {
	return NewInMemoryTechnologyRepository()
}

// Experience returns the experience repository for this transaction.
func (t *InMemoryTransaction) Experience() repository.ExperienceRepository {
	return NewInMemoryExperienceRepository()
}

// Service returns the service repository for this transaction.
func (t *InMemoryTransaction) Service() repository.ServiceRepository {
	return NewInMemoryServiceRepository()
}

// Commit finalizes the transaction changes.
func (t *InMemoryTransaction) Commit() error {
	return nil
}

// Rollback discards the transaction changes.
func (t *InMemoryTransaction) Rollback() error {
	return nil
}
