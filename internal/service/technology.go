// Package service provides business logic and application services for the portfolio website.
// It contains technology management services that handle CRUD operations, validation,
// and business rules for technology entities with filtering and categorization.
package service

import (
	"context"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// TechnologyService handles business logic for technologies.
type TechnologyService struct {
	repo         repository.TechnologyRepository
	validator    *CompoundValidator
	errorHandler *StandardServiceErrorHandlers
}

// NewTechnologyService creates a new technology service.
func NewTechnologyService(repo repository.TechnologyRepository) *TechnologyService {
	if repo == nil {
		panic("repository cannot be nil")
	}

	return &TechnologyService{
		repo:         repo,
		validator:    NewCompoundValidator(),
		errorHandler: NewStandardServiceErrorHandlers("TechnologyService"),
	}
}

// CreateTechnology creates a new technology with validation.
func (s *TechnologyService) CreateTechnology(ctx context.Context, name, category string, level domain.Level) (*domain.Technology, error) {
	// Validate and normalize input
	normalizedName, err := NewStringField("technology name", name).Required().MinLength(2).MaxLength(100).Validate()
	if err != nil {
		return nil, err
	}

	normalizedCategory, err := NewStringField("category", category).Required().MaxLength(50).Validate()
	if err != nil {
		return nil, err
	}

	// Check if technology already exists
	existing, err := s.repo.GetByName(ctx, normalizedName)
	if err == nil && existing != nil {
		return nil, s.errorHandler.Repository.HandleAlreadyExists("technology", normalizedName)
	}

	// Create new technology
	tech := domain.NewTechnology(normalizedName, normalizedCategory, level)
	tech.ID = generateID()

	if err := tech.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, tech); err != nil {
		return nil, s.errorHandler.Repository.HandleCreateError("technology", err)
	}

	return tech, nil
}

// GetTechnology retrieves a technology by ID.
func (s *TechnologyService) GetTechnology(ctx context.Context, id string) (*domain.Technology, error) {
	if err := s.validator.CommonRequestValidator.ValidateID(id, "technology"); err != nil {
		return nil, err
	}

	tech, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.errorHandler.Repository.HandleNotFound("technology", id)
	}

	return tech, nil
}

// ListTechnologies retrieves technologies with filtering and pagination.
func (s *TechnologyService) ListTechnologies(ctx context.Context, filter TechnologyFilter) ([]*domain.Technology, error) {
	// Validate list request parameters
	if err := s.validator.ValidateListRequest(filter.Limit, filter.Offset, filter.OrderBy, filter.OrderDir, "technology"); err != nil {
		return nil, err
	}

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
		return nil, s.errorHandler.Basic.HandleRepositoryListError("technologies", err)
	}

	return technologies, nil
}

// UpdateTechnology updates an existing technology.
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
		return nil, s.errorHandler.Repository.HandleUpdateError("technology", err)
	}

	return tech, nil
}

// DeleteTechnology removes a technology.
func (s *TechnologyService) DeleteTechnology(ctx context.Context, id string) error {
	if err := s.validator.CommonRequestValidator.ValidateID(id, "technology"); err != nil {
		return err
	}

	// Check if technology exists
	_, err := s.GetTechnology(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return s.errorHandler.Repository.HandleDeleteError("technology", err)
	}

	return nil
}

// GetTechnologiesByCategory retrieves technologies filtered by category.
func (s *TechnologyService) GetTechnologiesByCategory(ctx context.Context, category string) ([]*domain.Technology, error) {
	if err := s.validator.CommonRequestValidator.ValidateNonEmpty(category, "category"); err != nil {
		return nil, err
	}

	technologies, err := s.repo.GetByCategory(ctx, category)
	if err != nil {
		return nil, s.errorHandler.Repository.HandleCustomOperation("get technologies by category", err)
	}

	return technologies, nil
}

// GetTechnologiesByLevel retrieves technologies filtered by proficiency level.
func (s *TechnologyService) GetTechnologiesByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error) {
	if !level.IsValid() {
		return nil, domain.ErrInvalidInput("invalid technology level")
	}

	technologies, err := s.repo.GetByLevel(ctx, level)
	if err != nil {
		return nil, s.errorHandler.Repository.HandleCustomOperation("get technologies by level", err)
	}

	return technologies, nil
}

// TechnologyFilter represents filtering options for the service layer.
type TechnologyFilter struct {
	Category *string
	Level    *domain.Level
	Limit    *int
	Offset   *int
	OrderBy  *string
	OrderDir *string
}

// TechnologyUpdate represents fields that can be updated.
type TechnologyUpdate struct {
	Level       *domain.Level
	Description *string
	IconURL     *string
}
