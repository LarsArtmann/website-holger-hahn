// Package repository defines data access interfaces and contracts for the portfolio website.
// It provides abstractions for data persistence operations, transaction management,
// and repository patterns that can be implemented by various storage backends.
package repository

import (
	"context"
)

// Entity represents a generic entity interface
type Entity interface {
	GetID() string
}

// Repository defines common repository operations for any entity type
type Repository[T Entity] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity T) error

	// GetByID retrieves an entity by its ID
	GetByID(ctx context.Context, id string) (T, error)

	// Update updates an existing entity
	Update(ctx context.Context, entity T) error

	// Delete removes an entity by ID
	Delete(ctx context.Context, id string) error
}

// Filterable defines interface for entities that support filtering
type Filterable[T Entity, F any] interface {
	Repository[T]

	// List retrieves entities with optional filtering
	List(ctx context.Context, filter F) ([]T, error)
}

// UnitOfWork defines a unit of work pattern for transactional operations.
type UnitOfWork interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction.
type Transaction interface {
	// Technology returns the technology repository within this transaction
	Technology() TechnologyRepository

	// Experience returns the experience repository within this transaction
	Experience() ExperienceRepository

	// Service returns the service repository within this transaction
	Service() ServiceRepository

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}

// Repositories aggregates all repository interfaces.
type Repositories struct {
	Technology TechnologyRepository
	Experience ExperienceRepository
	Service    ServiceRepository
	UnitOfWork UnitOfWork
}

// HealthChecker defines interface for repository health checking.
type HealthChecker interface {
	// HealthCheck performs a health check on the repository
	HealthCheck(ctx context.Context) error
}

// Migrator defines interface for database migrations.
type Migrator interface {
	// Migrate runs database migrations
	Migrate(ctx context.Context) error

	// Rollback rolls back the last migration
	Rollback(ctx context.Context) error

	// Status returns the current migration status
	Status(ctx context.Context) (string, error)
}
