package infrastructure

import (
	"holger-hahn-website/internal/application"
	"holger-hahn-website/internal/domain"

	"github.com/samber/do"
)

// SetupContainer configures the dependency injection container.
func SetupContainer() *do.Injector {
	injector := do.New()

	// Register repositories
	do.Provide[domain.ContactRepository](injector, func(i *do.Injector) (domain.ContactRepository, error) {
		return NewMemoryContactRepository(), nil
	})

	// Register services
	do.Provide[domain.EmailService](injector, func(i *do.Injector) (domain.EmailService, error) {
		return NewSMTPEmailService(), nil
	})

	do.Provide[domain.LoggingService](injector, func(i *do.Injector) (domain.LoggingService, error) {
		return NewConsoleLoggingService("holger-hahn-website"), nil
	})

	// Register application services
	do.Provide[*application.ContactService](injector, func(i *do.Injector) (*application.ContactService, error) {
		contactRepo := do.MustInvoke[domain.ContactRepository](i)
		emailSvc := do.MustInvoke[domain.EmailService](i)
		logger := do.MustInvoke[domain.LoggingService](i)

		return application.NewContactService(contactRepo, emailSvc, logger), nil
	})

	return injector
}
