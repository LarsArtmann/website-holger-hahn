// Package service provides business logic and application services for the portfolio website.
// It contains experience management services that handle professional experience data,
// achievements, metrics tracking, and comprehensive career history operations.
package service

import (
	"context"
	"fmt"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// ExperienceService handles business logic for experiences.
type ExperienceService struct {
	repo         repository.ExperienceRepository
	techRepo     repository.TechnologyRepository
	validator    *CompoundValidator
	errorHandler *StandardServiceErrorHandlers
}

// NewExperienceService creates a new experience service.
func NewExperienceService(repo repository.ExperienceRepository, techRepo repository.TechnologyRepository) *ExperienceService {
	if repo == nil {
		panic("experience repository cannot be nil")
	}
	if techRepo == nil {
		panic("technology repository cannot be nil")
	}
	
	return &ExperienceService{
		repo:         repo,
		techRepo:     techRepo,
		validator:    NewCompoundValidator(),
		errorHandler: NewStandardServiceErrorHandlers("ExperienceService"),
	}
}

// CreateExperience creates a new experience with validation.
func (s *ExperienceService) CreateExperience(ctx context.Context, req CreateExperienceRequest) (*domain.Experience, error) {
	// Validate and normalize company name
	normalizedCompanyName, err := NewStringField("company name", req.CompanyName).Required().MinLength(2).MaxLength(100).Validate()
	if err != nil {
		return nil, err
	}

	// Validate and normalize position
	normalizedPosition, err := NewStringField("position", req.Position).Required().MinLength(2).MaxLength(100).Validate()
	if err != nil {
		return nil, err
	}

	// Validate and normalize optional fields
	normalizedLocation, err := NewStringField("location", req.Location).MaxLength(100).Validate()
	if err != nil {
		return nil, err
	}

	normalizedDescription, err := NewStringField("description", req.Description).MaxLength(2000).Validate()
	if err != nil {
		return nil, err
	}

	// Validate dates
	if err := s.validator.ValidateStartEndDates(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	// Create new experience
	experience := domain.NewExperience(
		normalizedCompanyName,
		normalizedPosition,
		normalizedDescription,
		normalizedLocation,
		req.StartDate,
		req.IsRemote,
	)
	experience.ID = generateID()

	// Set end date if provided
	if req.EndDate != nil {
		if err := experience.SetEndDate(*req.EndDate); err != nil {
			return nil, err
		}
	}

	if err := experience.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, experience); err != nil {
		return nil, s.errorHandler.Repository.HandleCreateError("experience", err)
	}

	return experience, nil
}

// GetExperience retrieves an experience by ID.
func (s *ExperienceService) GetExperience(ctx context.Context, id string) (*domain.Experience, error) {
	if err := s.validator.CommonRequestValidator.ValidateID(id, "experience"); err != nil {
		return nil, err
	}

	experience, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, s.errorHandler.Repository.HandleNotFound("experience", id)
	}

	return experience, nil
}

// ListExperiences retrieves experiences with filtering and pagination.
func (s *ExperienceService) ListExperiences(ctx context.Context, filter ExperienceFilter) ([]*domain.Experience, error) {
	repoFilter := repository.ExperienceFilter{
		CompanyName: filter.CompanyName,
		Position:    filter.Position,
		IsRemote:    filter.IsRemote,
		IsCurrent:   filter.IsCurrent,
		StartAfter:  filter.StartAfter,
		EndBefore:   filter.EndBefore,
		Technology:  filter.Technology,
		Location:    filter.Location,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
		OrderBy:     filter.OrderBy,
		OrderDir:    filter.OrderDir,
	}

	experiences, err := s.repo.List(ctx, repoFilter)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to list experiences: %v", err))
	}

	return experiences, nil
}

// GetCurrentExperiences retrieves all current experiences.
func (s *ExperienceService) GetCurrentExperiences(ctx context.Context) ([]*domain.Experience, error) {
	experiences, err := s.repo.GetCurrent(ctx)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to get current experiences: %v", err))
	}

	return experiences, nil
}

// AddTechnologyToExperience adds a technology to an experience.
func (s *ExperienceService) AddTechnologyToExperience(ctx context.Context, experienceID, technologyID string) error {
	// Get experience
	experience, err := s.GetExperience(ctx, experienceID)
	if err != nil {
		return err
	}

	// Get technology
	technology, err := s.techRepo.GetByID(ctx, technologyID)
	if err != nil {
		return domain.ErrNotFound("technology")
	}

	// Check if technology is already added
	for _, tech := range experience.Technologies {
		if tech.ID == technologyID {
			return domain.ErrConflict("technology already added to experience")
		}
	}

	// Add technology
	experience.AddTechnology(*technology)

	if err := s.repo.Update(ctx, experience); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to add technology to experience: %v", err))
	}

	return nil
}

// AddAchievementToExperience adds an achievement to an experience.
func (s *ExperienceService) AddAchievementToExperience(ctx context.Context, experienceID string, achievement AchievementRequest) error {
	// Get experience
	experience, err := s.GetExperience(ctx, experienceID)
	if err != nil {
		return err
	}

	// Create achievement
	var ach *domain.Achievement
	if achievement.Metrics != nil {
		ach = domain.NewAchievementWithMetrics(achievement.Title, achievement.Description, achievement.Impact, achievement.Metrics)
	} else {
		ach = domain.NewAchievement(achievement.Title, achievement.Description, achievement.Impact)
	}
	ach.ID = generateID()

	// Add achievement
	experience.AddAchievement(*ach)

	if err := s.repo.Update(ctx, experience); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to add achievement to experience: %v", err))
	}

	return nil
}

// EndExperience sets the end date for an experience.
func (s *ExperienceService) EndExperience(ctx context.Context, id string, endDate time.Time) error {
	experience, err := s.GetExperience(ctx, id)
	if err != nil {
		return err
	}

	if err := experience.SetEndDate(endDate); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, experience); err != nil {
		return domain.ErrInternal(fmt.Sprintf("failed to end experience: %v", err))
	}

	return nil
}

// GetExperiencesByCompany retrieves experiences for a specific company.
func (s *ExperienceService) GetExperiencesByCompany(ctx context.Context, companyName string) ([]*domain.Experience, error) {
	if companyName == "" {
		return nil, domain.ErrInvalidInput("company name cannot be empty")
	}

	experiences, err := s.repo.GetByCompany(ctx, companyName)
	if err != nil {
		return nil, domain.ErrInternal(fmt.Sprintf("failed to get experiences by company: %v", err))
	}

	return experiences, nil
}

// GetTotalExperience calculates total professional experience duration.
func (s *ExperienceService) GetTotalExperience(ctx context.Context) (time.Duration, error) {
	experiences, err := s.repo.List(ctx, repository.ExperienceFilter{})
	if err != nil {
		return 0, domain.ErrInternal(fmt.Sprintf("failed to calculate total experience: %v", err))
	}

	var total time.Duration
	for _, exp := range experiences {
		total += exp.Duration()
	}

	return total, nil
}

// CreateExperienceRequest represents a request to create an experience.
type CreateExperienceRequest struct {
	CompanyName string     `json:"company_name"`
	Position    string     `json:"position"`
	Description string     `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Location    string     `json:"location"`
	IsRemote    bool       `json:"is_remote"`
}

// ExperienceFilter represents filtering options for the service layer.
type ExperienceFilter struct {
	CompanyName *string
	Position    *string
	IsRemote    *bool
	IsCurrent   *bool
	StartAfter  *time.Time
	EndBefore   *time.Time
	Technology  *string
	Location    *string
	Limit       *int
	Offset      *int
	OrderBy     *string
	OrderDir    *string
}

// AchievementRequest represents a request to add an achievement.
type AchievementRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Impact      string          `json:"impact,omitempty"`
	Metrics     *domain.Metrics `json:"metrics,omitempty"`
}
