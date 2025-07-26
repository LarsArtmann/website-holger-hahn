// Package container provides in-memory repository implementations for development and testing.
// These implementations store data in memory and provide the same interface as production
// repositories for local development and unit testing purposes.
package container

import (
	"context"
	"fmt"
	"sync"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// InMemoryBaseCRUD provides common CRUD operations for in-memory repositories.
type InMemoryBaseCRUD[T repository.Entity] struct {
	entities map[string]T
	mu       sync.RWMutex
	typeName string
}

// NewInMemoryBaseCRUD creates a new base CRUD repository.
func NewInMemoryBaseCRUD[T repository.Entity](typeName string) *InMemoryBaseCRUD[T] {
	return &InMemoryBaseCRUD[T]{
		entities: make(map[string]T),
		typeName: typeName,
	}
}

// Create stores a new entity in the in-memory repository.
func (r *InMemoryBaseCRUD[T]) Create(ctx context.Context, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entities[entity.GetID()] = entity
	return nil
}

// GetByID retrieves an entity by its ID.
func (r *InMemoryBaseCRUD[T]) GetByID(ctx context.Context, id string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var zero T
	item, exists := r.entities[id]
	if !exists {
		return zero, domain.ErrNotFound(r.typeName)
	}

	return item, nil
}

// Update modifies an existing entity in the in-memory repository.
func (r *InMemoryBaseCRUD[T]) Update(ctx context.Context, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entities[entity.GetID()] = entity
	return nil
}

// Delete removes an entity from the in-memory repository by its ID.
func (r *InMemoryBaseCRUD[T]) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entities, id)
	return nil
}

// GetAll returns all entities in the repository (for filtering operations).
func (r *InMemoryBaseCRUD[T]) GetAll() map[string]T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Return a copy to prevent concurrent access issues
	result := make(map[string]T, len(r.entities))
	for k, v := range r.entities {
		result[k] = v
	}
	return result
}

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
	repo := &InMemoryExperienceRepository{
		exps: make(map[string]*domain.Experience),
	}
	repo.populateSampleData()
	return repo
}

// populateSampleData adds sample experiences with quantified metrics.
func (r *InMemoryExperienceRepository) populateSampleData() {
	// Fireblocks experience with comprehensive metrics
	fireblocksStart := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)
	fireblocksEnd := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	fireblocks := domain.NewExperience(
		"Fireblocks",
		"Senior Software Engineer",
		"Led development of digital asset custody solutions for institutional clients",
		"Tel Aviv, Israel",
		fireblocksStart,
		true,
	)
	fireblocks.ID = "fireblocks-001"
	// Error handling for SetEndDate: In data seeding context, date validation errors
	// indicate programming errors in test data, which should cause immediate failure
	if err := fireblocks.SetEndDate(fireblocksEnd); err != nil {
		panic(fmt.Sprintf("invalid fireblocks end date in test data: %v", err))
	}

	// Add achievements with quantified metrics
	testCoverageAchievement := domain.NewAchievementWithMetrics(
		"Improved Test Coverage by 25%",
		"Implemented comprehensive test automation strategy across critical financial systems",
		"Reduced production incidents by 60% and improved deployment confidence",
		&domain.Metrics{
			TestCoverage: domain.NewTestCoverageMetric(70.0, 95.0),
		},
	)
	testCoverageAchievement.ID = "achievement-001"

	deploymentAchievement := domain.NewAchievementWithMetrics(
		"Reduced Deployment Time by 87%",
		"Automated CI/CD pipeline with zero-downtime deployments and comprehensive testing",
		"Enabled daily deployments and faster time-to-market for critical features",
		&domain.Metrics{
			DeploymentTime: domain.NewDeploymentTimeMetric(120, 15), // 2h to 15min
		},
	)
	deploymentAchievement.ID = "achievement-002"

	systemReliabilityAchievement := domain.NewAchievementWithMetrics(
		"Achieved 99.9% System Uptime",
		"Implemented robust monitoring, alerting, and automated failover mechanisms",
		"Met enterprise SLA requirements and eliminated manual intervention incidents",
		&domain.Metrics{
			SystemReliability: domain.NewReliabilityMetric(99.9, 720, 5, 2), // 99.9%, 720h MTBF, 5min MTTR, 2 incidents/month
		},
	)
	systemReliabilityAchievement.ID = "achievement-003"

	productivityAchievement := domain.NewAchievementWithMetrics(
		"Increased Team Productivity by 40%",
		"Streamlined development processes and introduced modern DevOps practices",
		"Reduced lead time and improved team velocity while maintaining quality",
		&domain.Metrics{
			Productivity: domain.NewProductivityMetric(12, 48, 24, 40.0), // 12 deploys/week, 48h lead time, 24h cycle time, 40% efficiency gain
		},
	)
	productivityAchievement.ID = "achievement-004"

	costSavingsAchievement := domain.NewAchievementWithMetrics(
		"Generated $2.4M Annual Cost Savings",
		"Optimized infrastructure costs and eliminated operational inefficiencies",
		"Achieved 300% ROI through automation and process improvements",
		&domain.Metrics{
			CostSavings: domain.NewCostMetric(200000, 2400000, 300.0, 4), // $200k monthly, $2.4M annual, 300% ROI, 4 months payback
		},
	)
	costSavingsAchievement.ID = "achievement-005"

	fireblocks.AddAchievement(*testCoverageAchievement)
	fireblocks.AddAchievement(*deploymentAchievement)
	fireblocks.AddAchievement(*systemReliabilityAchievement)
	fireblocks.AddAchievement(*productivityAchievement)
	fireblocks.AddAchievement(*costSavingsAchievement)

	r.exps[fireblocks.ID] = fireblocks

	// Current Consulting experience with metrics
	consultingStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	consulting := domain.NewExperience(
		"Independent Consultant",
		"Senior Software Engineer & Technical Consultant",
		"Providing specialized expertise in digital asset infrastructure and secure financial systems",
		"Remote",
		consultingStart,
		true,
	)
	consulting.ID = "consulting-001"

	consultingAchievement := domain.NewAchievementWithMetrics(
		"Accelerated Client Deployment by 75%",
		"Implemented standardized development practices and automated testing frameworks",
		"Enabled faster time-to-market and reduced technical risk for financial product launches",
		&domain.Metrics{
			DeploymentTime: domain.NewDeploymentTimeMetric(480, 120), // 8h to 2h
			TestCoverage:   domain.NewTestCoverageMetric(60.0, 85.0),
			Productivity:   domain.NewProductivityMetric(8, 72, 36, 35.0), // 8 deploys/week, 72h lead time, 36h cycle time, 35% improvement
		},
	)
	consultingAchievement.ID = "achievement-006"

	consulting.AddAchievement(*consultingAchievement)
	r.exps[consulting.ID] = consulting

	// Previous Banking experience with legacy system metrics
	bankingStart := time.Date(2018, 3, 1, 0, 0, 0, 0, time.UTC)
	bankingEnd := time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC)
	banking := domain.NewExperience(
		"Major European Bank",
		"Software Engineer",
		"Developed and maintained critical banking infrastructure and customer-facing applications",
		"Frankfurt, Germany",
		bankingStart,
		false,
	)
	banking.ID = "banking-001"
	// Error handling for SetEndDate: In data seeding context, date validation errors
	// indicate programming errors in test data, which should cause immediate failure
	if err := banking.SetEndDate(bankingEnd); err != nil {
		panic(fmt.Sprintf("invalid banking end date in test data: %v", err))
	}

	legacyModernizationAchievement := domain.NewAchievementWithMetrics(
		"Modernized Legacy Systems with 50% Performance Gain",
		"Refactored critical payment processing systems to improve reliability and performance",
		"Reduced transaction processing time and improved customer experience",
		&domain.Metrics{
			SystemReliability: domain.NewReliabilityMetric(99.5, 480, 15, 5),   // 99.5%, 480h MTBF, 15min MTTR, 5 incidents/month
			CostSavings:       domain.NewCostMetric(150000, 1800000, 250.0, 6), // $150k monthly, $1.8M annual, 250% ROI, 6 months payback
		},
	)
	legacyModernizationAchievement.ID = "achievement-007"

	banking.AddAchievement(*legacyModernizationAchievement)
	r.exps[banking.ID] = banking
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
