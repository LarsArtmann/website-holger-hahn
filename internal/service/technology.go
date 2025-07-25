package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// TechnologyService handles business logic for technologies
type TechnologyService struct {
	repo repository.TechnologyRepository
}

// NewTechnologyService creates a new technology service
func NewTechnologyService(repo repository.TechnologyRepository) *TechnologyService {
	return &TechnologyService{
		repo: repo,
	}
}

// CreateTechnology creates a new technology with validation
func (s *TechnologyService) CreateTechnology(ctx context.Context, name, category string, level domain.Level) (*domain.Technology, error) {
	// Normalize input
	name = strings.TrimSpace(name)
	category = strings.TrimSpace(category)
	
	if name == "" {
		return nil, domain.ErrInvalidInput("technology name cannot be empty")
	}
	
	// Check if technology already exists
	existing, err := s.repo.GetByName(ctx, name)
	if err == nil && existing != nil {
		return nil, domain.ErrConflict(fmt.Sprintf("technology '%s' already exists", name))
	}
	
	// Create new technology
	tech := domain.NewTechnology(name, category, level)
	tech.ID = generateID()
	
	if err := tech.Validate(); err != nil {
		return nil, err
	}
	
	if err := s.repo.Create(ctx, tech); err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to create technology: %v", err))
	}
	
	return tech, nil
}

// GetTechnology retrieves a technology by ID
func (s *TechnologyService) GetTechnology(ctx context.Context, id string) (*domain.Technology, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput("technology ID cannot be empty")
	}
	
	tech, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound("technology")
	}
	
	return tech, nil
}

// ListTechnologies retrieves technologies with filtering and pagination
func (s *TechnologyService) ListTechnologies(ctx context.Context, filter TechnologyFilter) ([]*domain.Technology, error) {
	repoFilter := repository.TechnologyFilter{
		Category: filter.Category,
		Level:    filter.Level,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
		OrderBy:  filter.OrderBy,
		OrderDir: filter.OrderDir,
	}
	
	technologies, err := s.repo.List(ctx, repoFilter)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to list technologies: %v", err))
	}
	
	return technologies, nil
}

// UpdateTechnology updates an existing technology
func (s *TechnologyService) UpdateTechnology(ctx context.Context, id string, updates TechnologyUpdate) (*domain.Technology, error) {
	tech, err := s.GetTechnology(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Apply updates
	if updates.Level != nil {
		if err := tech.UpdateLevel(*updates.Level); err != nil {
			return nil, err
		}
	}
	
	if updates.Description != nil {
		tech.Description = *updates.Description
		tech.UpdatedAt = time.Now()
	}
	
	if updates.IconURL != nil {
		tech.IconURL = *updates.IconURL
		tech.UpdatedAt = time.Now()
	}
	
	if err := tech.Validate(); err != nil {
		return nil, err
	}
	
	if err := s.repo.Update(ctx, tech); err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to update technology: %v", err))
	}
	
	return tech, nil
}

// DeleteTechnology removes a technology
func (s *TechnologyService) DeleteTechnology(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput("technology ID cannot be empty")
	}
	
	// Check if technology exists
	_, err := s.GetTechnology(ctx, id)
	if err != nil {
		return err
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to delete technology: %v", err))
	}
	
	return nil
}

// GetTechnologiesByCategory retrieves technologies filtered by category
func (s *TechnologyService) GetTechnologiesByCategory(ctx context.Context, category string) ([]*domain.Technology, error) {
	if category == "" {
		return nil, domain.ErrInvalidInput("category cannot be empty")
	}
	
	technologies, err := s.repo.GetByCategory(ctx, category)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to get technologies by category: %v", err))
	}
	
	return technologies, nil
}

// GetTechnologiesByLevel retrieves technologies filtered by proficiency level
func (s *TechnologyService) GetTechnologiesByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error) {
	if !level.IsValid() {
		return nil, domain.ErrInvalidInput("invalid technology level")
	}
	
	technologies, err := s.repo.GetByLevel(ctx, level)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to get technologies by level: %v", err))
	}
	
	return technologies, nil
}

// TechnologyFilter represents filtering options for the service layer
type TechnologyFilter struct {
	Category *string
	Level    *domain.Level
	Limit    *int
	Offset   *int
	OrderBy  *string
	OrderDir *string
}

// TechnologyUpdate represents fields that can be updated
type TechnologyUpdate struct {
	Level       *domain.Level
	Description *string
	IconURL     *string
}