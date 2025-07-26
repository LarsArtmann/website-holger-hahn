package testutil

import (
	"context"
	"fmt"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// MockTechnologyRepository is a mock implementation for testing
type MockTechnologyRepository struct {
	technologies map[string]*domain.Technology
	callLog      []string
	errorMode    bool
	errorMsg     string
}

// NewMockTechnologyRepository creates a new mock technology repository
func NewMockTechnologyRepository() *MockTechnologyRepository {
	return &MockTechnologyRepository{
		technologies: make(map[string]*domain.Technology),
		callLog:      make([]string, 0),
	}
}

// SetErrorMode enables error mode for testing error scenarios
func (m *MockTechnologyRepository) SetErrorMode(enabled bool, errorMsg string) {
	m.errorMode = enabled
	m.errorMsg = errorMsg
}

// GetCallLog returns the log of method calls
func (m *MockTechnologyRepository) GetCallLog() []string {
	return m.callLog
}

// ClearCallLog clears the call log
func (m *MockTechnologyRepository) ClearCallLog() {
	m.callLog = make([]string, 0)
}

// PreloadTechnologies preloads the mock with test data
func (m *MockTechnologyRepository) PreloadTechnologies(technologies []*domain.Technology) {
	for _, tech := range technologies {
		m.technologies[tech.ID] = tech
	}
}

func (m *MockTechnologyRepository) Create(ctx context.Context, tech *domain.Technology) error {
	m.callLog = append(m.callLog, fmt.Sprintf("Create(%s)", tech.ID))
	if m.errorMode {
		return fmt.Errorf(m.errorMsg)
	}
	m.technologies[tech.ID] = tech
	return nil
}

func (m *MockTechnologyRepository) GetByID(ctx context.Context, id string) (*domain.Technology, error) {
	m.callLog = append(m.callLog, fmt.Sprintf("GetByID(%s)", id))
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}
	tech, exists := m.technologies[id]
	if !exists {
		return nil, domain.ErrNotFound("technology")
	}
	return tech, nil
}

func (m *MockTechnologyRepository) GetByName(ctx context.Context, name string) (*domain.Technology, error) {
	m.callLog = append(m.callLog, fmt.Sprintf("GetByName(%s)", name))
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}
	for _, tech := range m.technologies {
		if tech.Name == name {
			return tech, nil
		}
	}
	return nil, domain.ErrNotFound("technology")
}

func (m *MockTechnologyRepository) List(ctx context.Context, filter repository.TechnologyFilter) ([]*domain.Technology, error) {
	m.callLog = append(m.callLog, "List(filter)")
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}

	var result []*domain.Technology
	for _, tech := range m.technologies {
		// Apply filters
		if filter.Category != nil && tech.Category != *filter.Category {
			continue
		}
		if filter.Level != nil && tech.Level != *filter.Level {
			continue
		}
		result = append(result, tech)
	}

	// Apply pagination
	if filter.Offset != nil && filter.Limit != nil {
		offset := *filter.Offset
		limit := *filter.Limit
		if offset >= len(result) {
			return []*domain.Technology{}, nil
		}
		end := offset + limit
		if end > len(result) {
			end = len(result)
		}
		result = result[offset:end]
	}

	return result, nil
}

func (m *MockTechnologyRepository) Update(ctx context.Context, tech *domain.Technology) error {
	m.callLog = append(m.callLog, fmt.Sprintf("Update(%s)", tech.ID))
	if m.errorMode {
		return fmt.Errorf(m.errorMsg)
	}
	m.technologies[tech.ID] = tech
	return nil
}

func (m *MockTechnologyRepository) Delete(ctx context.Context, id string) error {
	m.callLog = append(m.callLog, fmt.Sprintf("Delete(%s)", id))
	if m.errorMode {
		return fmt.Errorf(m.errorMsg)
	}
	delete(m.technologies, id)
	return nil
}

func (m *MockTechnologyRepository) GetByCategory(ctx context.Context, category string) ([]*domain.Technology, error) {
	m.callLog = append(m.callLog, fmt.Sprintf("GetByCategory(%s)", category))
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}

	var result []*domain.Technology
	for _, tech := range m.technologies {
		if tech.Category == category {
			result = append(result, tech)
		}
	}
	return result, nil
}

func (m *MockTechnologyRepository) GetByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error) {
	m.callLog = append(m.callLog, fmt.Sprintf("GetByLevel(%s)", level))
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}

	var result []*domain.Technology
	for _, tech := range m.technologies {
		if tech.Level == level {
			result = append(result, tech)
		}
	}
	return result, nil
}

// MockExperienceRepository is a mock implementation for testing
type MockExperienceRepository struct {
	experiences map[string]*domain.Experience
	callLog     []string
	errorMode   bool
	errorMsg    string
}

// NewMockExperienceRepository creates a new mock experience repository
func NewMockExperienceRepository() *MockExperienceRepository {
	return &MockExperienceRepository{
		experiences: make(map[string]*domain.Experience),
		callLog:     make([]string, 0),
	}
}

// SetErrorMode enables error mode for testing error scenarios
func (m *MockExperienceRepository) SetErrorMode(enabled bool, errorMsg string) {
	m.errorMode = enabled
	m.errorMsg = errorMsg
}

// GetCallLog returns the log of method calls
func (m *MockExperienceRepository) GetCallLog() []string {
	return m.callLog
}

// ClearCallLog clears the call log
func (m *MockExperienceRepository) ClearCallLog() {
	m.callLog = make([]string, 0)
}

// PreloadExperiences preloads the mock with test data
func (m *MockExperienceRepository) PreloadExperiences(experiences []*domain.Experience) {
	for _, exp := range experiences {
		m.experiences[exp.ID] = exp
	}
}

func (m *MockExperienceRepository) Create(ctx context.Context, exp *domain.Experience) error {
	m.callLog = append(m.callLog, fmt.Sprintf("Create(%s)", exp.ID))
	if m.errorMode {
		return fmt.Errorf(m.errorMsg)
	}
	m.experiences[exp.ID] = exp
	return nil
}

func (m *MockExperienceRepository) GetByID(ctx context.Context, id string) (*domain.Experience, error) {
	m.callLog = append(m.callLog, fmt.Sprintf("GetByID(%s)", id))
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}
	exp, exists := m.experiences[id]
	if !exists {
		return nil, domain.ErrNotFound("experience")
	}
	return exp, nil
}

func (m *MockExperienceRepository) List(ctx context.Context, filter repository.ExperienceFilter) ([]*domain.Experience, error) {
	m.callLog = append(m.callLog, "List(filter)")
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}

	var result []*domain.Experience
	for _, exp := range m.experiences {
		// Apply filters
		if filter.CompanyName != nil && exp.CompanyName != *filter.CompanyName {
			continue
		}
		if filter.Position != nil && exp.Position != *filter.Position {
			continue
		}
		if filter.IsRemote != nil && exp.IsRemote != *filter.IsRemote {
			continue
		}
		if filter.IsCurrent != nil && exp.IsCurrent() != *filter.IsCurrent {
			continue
		}
		result = append(result, exp)
	}

	// Apply pagination
	if filter.Offset != nil && filter.Limit != nil {
		offset := *filter.Offset
		limit := *filter.Limit
		if offset >= len(result) {
			return []*domain.Experience{}, nil
		}
		end := offset + limit
		if end > len(result) {
			end = len(result)
		}
		result = result[offset:end]
	}

	return result, nil
}

func (m *MockExperienceRepository) Update(ctx context.Context, exp *domain.Experience) error {
	m.callLog = append(m.callLog, fmt.Sprintf("Update(%s)", exp.ID))
	if m.errorMode {
		return fmt.Errorf(m.errorMsg)
	}
	m.experiences[exp.ID] = exp
	return nil
}

func (m *MockExperienceRepository) Delete(ctx context.Context, id string) error {
	m.callLog = append(m.callLog, fmt.Sprintf("Delete(%s)", id))
	if m.errorMode {
		return fmt.Errorf(m.errorMsg)
	}
	delete(m.experiences, id)
	return nil
}

func (m *MockExperienceRepository) GetCurrent(ctx context.Context) ([]*domain.Experience, error) {
	m.callLog = append(m.callLog, "GetCurrent()")
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}

	var result []*domain.Experience
	for _, exp := range m.experiences {
		if exp.IsCurrent() {
			result = append(result, exp)
		}
	}
	return result, nil
}

func (m *MockExperienceRepository) GetByCompany(ctx context.Context, companyName string) ([]*domain.Experience, error) {
	m.callLog = append(m.callLog, fmt.Sprintf("GetByCompany(%s)", companyName))
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}

	var result []*domain.Experience
	for _, exp := range m.experiences {
		if exp.CompanyName == companyName {
			result = append(result, exp)
		}
	}
	return result, nil
}

// MockServiceRepository is a mock implementation for testing
type MockServiceRepository struct {
	services  map[string]interface{} // Generic services for portfolio service
	callLog   []string
	errorMode bool
	errorMsg  string
}

// NewMockServiceRepository creates a new mock service repository
func NewMockServiceRepository() *MockServiceRepository {
	return &MockServiceRepository{
		services: make(map[string]interface{}),
		callLog:  make([]string, 0),
	}
}

// SetErrorMode enables error mode for testing error scenarios
func (m *MockServiceRepository) SetErrorMode(enabled bool, errorMsg string) {
	m.errorMode = enabled
	m.errorMsg = errorMsg
}

// GetCallLog returns the log of method calls
func (m *MockServiceRepository) GetCallLog() []string {
	return m.callLog
}

// ClearCallLog clears the call log
func (m *MockServiceRepository) ClearCallLog() {
	m.callLog = make([]string, 0)
}

// GetActiveServices returns mock active services
func (m *MockServiceRepository) GetActiveServices(ctx context.Context) (interface{}, error) {
	m.callLog = append(m.callLog, "GetActiveServices()")
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}
	return []string{"Blockchain Consulting", "Smart Contract Development", "DeFi Solutions"}, nil
}

// ListServices returns mock services list
func (m *MockServiceRepository) ListServices(ctx context.Context, filter interface{}) (interface{}, error) {
	m.callLog = append(m.callLog, "ListServices(filter)")
	if m.errorMode {
		return nil, fmt.Errorf(m.errorMsg)
	}
	return []string{"Blockchain Consulting", "Smart Contract Development", "DeFi Solutions", "Technical Auditing"}, nil
}
