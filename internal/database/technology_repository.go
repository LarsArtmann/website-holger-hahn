// Package database provides database repository implementations using sqlc generated code.
// This file implements the TechnologyRepository interface with SQLite backend.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/repository"
)

// TechnologyRepository implements repository.TechnologyRepository using sqlc generated code.
type TechnologyRepository struct {
	queries *Queries
}

// NewTechnologyRepository creates a new database technology repository.
func NewTechnologyRepository(queries *Queries) *TechnologyRepository {
	return &TechnologyRepository{
		queries: queries,
	}
}

// Create creates a new technology.
func (r *TechnologyRepository) Create(ctx context.Context, entity *domain.Technology) error {
	if entity == nil {
		return fmt.Errorf("technology cannot be nil")
	}

	params := CreateTechnologyParams{
		Name:             entity.Name,
		Category:         entity.Category,
		ProficiencyLevel: string(entity.Level),
		IconClass:        nullStringFromString(entity.IconURL), // Map IconURL to IconClass
		ColorScheme:      sql.NullString{},                     // Default empty
		Description:      nullStringFromString(entity.Description),
		SortOrder:        sql.NullInt64{Int64: 0, Valid: true}, // Default sort order
	}

	created, err := r.queries.CreateTechnology(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create technology: %w", err)
	}

	// Update entity with generated values
	entity.ID = created.ID
	if created.CreatedAt.Valid {
		entity.CreatedAt = created.CreatedAt.Time
	}
	if created.UpdatedAt.Valid {
		entity.UpdatedAt = created.UpdatedAt.Time
	}

	return nil
}

// GetByID retrieves a technology by its ID.
func (r *TechnologyRepository) GetByID(ctx context.Context, id string) (*domain.Technology, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	dbTech, err := r.queries.GetTechnology(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("technology not found")
		}
		return nil, fmt.Errorf("failed to get technology: %w", err)
	}

	return r.toDomainTechnology(dbTech), nil
}

// Update updates an existing technology.
func (r *TechnologyRepository) Update(ctx context.Context, entity *domain.Technology) error {
	if entity == nil {
		return fmt.Errorf("technology cannot be nil")
	}

	params := UpdateTechnologyParams{
		Name:             entity.Name,
		Category:         entity.Category,
		ProficiencyLevel: string(entity.Level),
		IconClass:        nullStringFromString(entity.IconURL),
		ColorScheme:      sql.NullString{},
		Description:      nullStringFromString(entity.Description),
		SortOrder:        sql.NullInt64{Int64: 0, Valid: true},
		ID:               entity.ID,
	}

	updated, err := r.queries.UpdateTechnology(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update technology: %w", err)
	}

	// Update timestamps
	if updated.UpdatedAt.Valid {
		entity.UpdatedAt = updated.UpdatedAt.Time
	}

	return nil
}

// Delete removes a technology by ID.
func (r *TechnologyRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if err := r.queries.DeleteTechnology(ctx, id); err != nil {
		return fmt.Errorf("failed to delete technology: %w", err)
	}

	return nil
}

// List retrieves technologies with optional filtering.
func (r *TechnologyRepository) List(ctx context.Context, filter repository.TechnologyFilter) ([]*domain.Technology, error) {
	var dbTechs []Technology
	var err error

	// Apply filters based on what's provided
	if filter.Category != nil {
		dbTechs, err = r.queries.ListTechnologiesByCategory(ctx, *filter.Category)
	} else if filter.Level != nil {
		dbTechs, err = r.queries.ListTechnologiesByLevel(ctx, string(*filter.Level))
	} else {
		dbTechs, err = r.queries.ListTechnologies(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list technologies: %w", err)
	}

	// Convert to domain entities
	technologies := make([]*domain.Technology, len(dbTechs))
	for i, dbTech := range dbTechs {
		technologies[i] = r.toDomainTechnology(dbTech)
	}

	// Apply client-side pagination if needed (since sqlc doesn't support LIMIT/OFFSET in all queries)
	if filter.Offset != nil && *filter.Offset > 0 {
		start := *filter.Offset
		if start >= len(technologies) {
			return []*domain.Technology{}, nil
		}
		technologies = technologies[start:]
	}

	if filter.Limit != nil && *filter.Limit > 0 {
		end := *filter.Limit
		if end > len(technologies) {
			end = len(technologies)
		}
		technologies = technologies[:end]
	}

	return technologies, nil
}

// GetByName retrieves a technology by its name.
func (r *TechnologyRepository) GetByName(ctx context.Context, name string) (*domain.Technology, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	dbTech, err := r.queries.GetTechnologyByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("technology not found")
		}
		return nil, fmt.Errorf("failed to get technology by name: %w", err)
	}

	return r.toDomainTechnology(dbTech), nil
}

// GetByCategory retrieves technologies by category.
func (r *TechnologyRepository) GetByCategory(ctx context.Context, category string) ([]*domain.Technology, error) {
	dbTechs, err := r.queries.ListTechnologiesByCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to get technologies by category: %w", err)
	}

	technologies := make([]*domain.Technology, len(dbTechs))
	for i, dbTech := range dbTechs {
		technologies[i] = r.toDomainTechnology(dbTech)
	}

	return technologies, nil
}

// GetByLevel retrieves technologies by proficiency level.
func (r *TechnologyRepository) GetByLevel(ctx context.Context, level domain.Level) ([]*domain.Technology, error) {
	dbTechs, err := r.queries.ListTechnologiesByLevel(ctx, string(level))
	if err != nil {
		return nil, fmt.Errorf("failed to get technologies by level: %w", err)
	}

	technologies := make([]*domain.Technology, len(dbTechs))
	for i, dbTech := range dbTechs {
		technologies[i] = r.toDomainTechnology(dbTech)
	}

	return technologies, nil
}

// GetID implements repository.Entity interface.
func (t *Technology) GetID() string {
	return t.ID
}

// toDomainTechnology converts a database Technology to a domain Technology.
func (r *TechnologyRepository) toDomainTechnology(dbTech Technology) *domain.Technology {
	tech := &domain.Technology{
		ID:       dbTech.ID,
		Name:     dbTech.Name,
		Category: dbTech.Category,
		Level:    domain.Level(dbTech.ProficiencyLevel),
	}

	if dbTech.IconClass.Valid {
		tech.IconURL = dbTech.IconClass.String
	}

	if dbTech.Description.Valid {
		tech.Description = dbTech.Description.String
	}

	if dbTech.CreatedAt.Valid {
		tech.CreatedAt = dbTech.CreatedAt.Time
	} else {
		tech.CreatedAt = time.Now().UTC()
	}

	if dbTech.UpdatedAt.Valid {
		tech.UpdatedAt = dbTech.UpdatedAt.Time
	} else {
		tech.UpdatedAt = time.Now().UTC()
	}

	return tech
}
