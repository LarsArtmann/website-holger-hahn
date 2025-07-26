// Package container provides dependency injection setup for the portfolio website.
// It configures and manages the dependency injection container with all required
// services, repositories, and infrastructure components.
package container

import (
	"github.com/samber/do"

	"holger-hahn-website/internal/config"
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

	// Register repository implementations.
	// Note: In a real application, you would register actual implementations.
	// For now, we're registering placeholder providers that return nil.
	// These would be replaced with actual database implementations.

	do.Provide(c.injector, func(_ *do.Injector) (repository.TechnologyRepository, error) {
		return NewInMemoryTechnologyRepository(), nil
	})

	do.Provide(c.injector, func(_ *do.Injector) (repository.ExperienceRepository, error) {
		return NewInMemoryExperienceRepository(), nil
	})

	do.Provide(c.injector, func(_ *do.Injector) (repository.ServiceRepository, error) {
		return NewInMemoryServiceRepository(), nil
	})

	do.Provide(c.injector, func(_ *do.Injector) (repository.UnitOfWork, error) {
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
