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

func (r *InMemoryTechnologyRepository) Create(ctx context.Context, technology *domain.Technology) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.techs[technology.ID] = technology

	return nil
}

func (r *InMemoryTechnologyRepository) GetByID(ctx context.Context, id string) (*domain.Technology, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tech, exists := r.techs[id]
	if !exists {
		return nil, domain.ErrNotFound("technology")
	}

	return tech, nil
}

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

func (r *InMemoryTechnologyRepository) Update(ctx context.Context, technology *domain.Technology) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.techs[technology.ID] = technology

	return nil
}

func (r *InMemoryTechnologyRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.techs, id)

	return nil
}

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

func (r *InMemoryExperienceRepository) Create(ctx context.Context, experience *domain.Experience) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exps[experience.ID] = experience

	return nil
}

func (r *InMemoryExperienceRepository) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exp, exists := r.exps[id]
	if !exists {
		return nil, domain.ErrNotFound("experience")
	}

	return exp, nil
}

func (r *InMemoryExperienceRepository) List(ctx context.Context, filter repository.ExperienceFilter) ([]*domain.Experience, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Experience
	for _, exp := range r.exps {
		result = append(result, exp)
	}

	return result, nil
}

func (r *InMemoryExperienceRepository) Update(ctx context.Context, experience *domain.Experience) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exps[experience.ID] = experience

	return nil
}

func (r *InMemoryExperienceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.exps, id)

	return nil
}

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

func (r *InMemoryServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[service.ID] = service

	return nil
}

func (r *InMemoryServiceRepository) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc, exists := r.services[id]
	if !exists {
		return nil, domain.ErrNotFound("service")
	}

	return svc, nil
}

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

func (r *InMemoryServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.services[service.ID] = service

	return nil
}

func (r *InMemoryServiceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, id)

	return nil
}

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

func (u *InMemoryUnitOfWork) Begin(ctx context.Context) (repository.Transaction, error) {
	return &InMemoryTransaction{}, nil
}

// InMemoryTransaction is a simple in-memory transaction implementation.
type InMemoryTransaction struct{}

func (t *InMemoryTransaction) Technology() repository.TechnologyRepository {
	return NewInMemoryTechnologyRepository()
}

func (t *InMemoryTransaction) Experience() repository.ExperienceRepository {
	return NewInMemoryExperienceRepository()
}

func (t *InMemoryTransaction) Service() repository.ServiceRepository {
	return NewInMemoryServiceRepository()
}

func (t *InMemoryTransaction) Commit() error {
	return nil
}

func (t *InMemoryTransaction) Rollback() error {
	return nil
}
