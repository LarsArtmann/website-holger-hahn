// Package container provides dependency injection setup for the portfolio website.
// It configures and manages the dependency injection container with all required
// services, repositories, and infrastructure components.
package container

import (
	"context"
	"log"

	"github.com/samber/do"
	"holger-hahn-website/internal/application"
	"holger-hahn-website/internal/config"
	"holger-hahn-website/internal/database"
	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/infrastructure"
	"holger-hahn-website/internal/repository"
	"holger-hahn-website/internal/service"
)

// Container wraps the DI container with our application-specific setup.
type Container struct {
	injector *do.Injector
}

// New creates a new container with all dependencies registered.
func New() *Container {
	injector := do.New()

	c := &Container{
		injector: injector,
	}

	c.registerDependencies()

	return c
}

// GetInjector returns the underlying DI injector.
func (c *Container) GetInjector() *do.Injector {
	return c.injector
}

// registerDependencies registers all application dependencies.
func (c *Container) registerDependencies() {
	// Register configuration.
	do.Provide(c.injector, func(_ *do.Injector) (*config.Config, error) {
		cfg := config.LoadConfig()
		if err := cfg.Validate(); err != nil {
			return nil, err
		}

		return cfg, nil
	})

	// Register database manager and initialize database
	do.Provide(c.injector, func(_ *do.Injector) (*database.DatabaseManager, error) {
		dbConfig := database.DefaultConfig()
		dbManager, err := database.NewDatabaseManager(dbConfig)
		if err != nil {
			return nil, err
		}

		// Run migrations on startup
		ctx := context.Background()
		if err := dbManager.Migrate(ctx); err != nil {
			log.Printf("Warning: Database migration failed: %v", err)
			// Don't fail startup, but log the warning
		}

		return dbManager, nil
	})

	// Register database repositories (replacing in-memory implementations)
	do.Provide(c.injector, func(i *do.Injector) (repository.TechnologyRepository, error) {
		dbManager := do.MustInvoke[*database.DatabaseManager](i)
		return database.NewTechnologyRepository(dbManager.Queries()), nil
	})

	do.Provide(c.injector, func(_ *do.Injector) (repository.ExperienceRepository, error) {
		// TODO: Implement database experience repository when available
		return NewInMemoryExperienceRepository(), nil
	})

	do.Provide(c.injector, func(_ *do.Injector) (repository.ServiceRepository, error) {
		// TODO: Implement database service repository when available  
		return NewInMemoryServiceRepository(), nil
	})

	do.Provide(c.injector, func(_ *do.Injector) (repository.UnitOfWork, error) {
		// TODO: Implement database unit of work when available
		return NewInMemoryUnitOfWork(), nil
	})

	// Register services.
	do.Provide(c.injector, func(i *do.Injector) (*service.TechnologyService, error) {
		techRepo := do.MustInvoke[repository.TechnologyRepository](i)
		return service.NewTechnologyService(techRepo), nil
	})

	do.Provide(c.injector, func(i *do.Injector) (*service.ExperienceService, error) {
		expRepo := do.MustInvoke[repository.ExperienceRepository](i)
		techRepo := do.MustInvoke[repository.TechnologyRepository](i)

		return service.NewExperienceService(expRepo, techRepo), nil
	})

	do.Provide(c.injector, func(i *do.Injector) (*service.PortfolioService, error) {
		serviceRepo := do.MustInvoke[repository.ServiceRepository](i)
		techRepo := do.MustInvoke[repository.TechnologyRepository](i)

		return service.NewPortfolioService(serviceRepo, techRepo), nil
	})

	// Register aggregated repositories struct.
	do.Provide(c.injector, func(i *do.Injector) (*repository.Repositories, error) {
		return &repository.Repositories{
			Technology: do.MustInvoke[repository.TechnologyRepository](i),
			Experience: do.MustInvoke[repository.ExperienceRepository](i),
			Service:    do.MustInvoke[repository.ServiceRepository](i),
			UnitOfWork: do.MustInvoke[repository.UnitOfWork](i),
		}, nil
	})

	// Register contact-related dependencies
	// Contact repository
	do.Provide(c.injector, func(_ *do.Injector) (domain.ContactRepository, error) {
		return infrastructure.NewMemoryContactRepository(), nil
	})

	// Email service
	do.Provide(c.injector, func(_ *do.Injector) (domain.EmailService, error) {
		return infrastructure.NewSMTPEmailService(), nil
	})

	// Logging service
	do.Provide(c.injector, func(_ *do.Injector) (domain.LoggingService, error) {
		return infrastructure.NewConsoleLoggingService("holger-hahn-website"), nil
	})

	// Contact application service
	do.Provide(c.injector, func(i *do.Injector) (*application.ContactService, error) {
		contactRepo := do.MustInvoke[domain.ContactRepository](i)
		emailSvc := do.MustInvoke[domain.EmailService](i)
		logger := do.MustInvoke[domain.LoggingService](i)

		return application.NewContactService(contactRepo, emailSvc, logger), nil
	})
}

// Shutdown gracefully shuts down the container and cleans up resources.
func (c *Container) Shutdown() error {
	return c.injector.Shutdown()
}

// MustGet retrieves a dependency from the container and panics if it fails.
func MustGet[T any](c *Container) T {
	return do.MustInvoke[T](c.injector)
}

// Get retrieves a dependency from the container.
func Get[T any](c *Container) (T, error) {
	return do.Invoke[T](c.injector)
}
