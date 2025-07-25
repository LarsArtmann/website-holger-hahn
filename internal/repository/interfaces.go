package repository

import (
	"context"
)

// UnitOfWork defines a unit of work pattern for transactional operations
type UnitOfWork interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction
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

// Repositories aggregates all repository interfaces
type Repositories struct {
	Technology TechnologyRepository
	Experience ExperienceRepository
	Service    ServiceRepository
	UnitOfWork UnitOfWork
}

// HealthChecker defines interface for repository health checking
type HealthChecker interface {
	// HealthCheck performs a health check on the repository
	HealthCheck(ctx context.Context) error
}

// Migrator defines interface for database migrations
type Migrator interface {
	// Migrate runs database migrations
	Migrate(ctx context.Context) error
	
	// Rollback rolls back the last migration
	Rollback(ctx context.Context) error
	
	// Status returns the current migration status
	Status(ctx context.Context) (string, error)
}