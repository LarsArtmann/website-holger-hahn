# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Professional website for Holger M. Hahn (Solution Architect & Test Engineer) built with Go + Templ + HTMX stack. The site showcases expertise in Digital Assets, Tokenization, and Web3 technologies, targeting freelance projects from financial institutions.

## Development Commands

### Setup & Development
```bash
just setup                    # Install all dependencies (Bun + Go)
just dev                     # Start development server with hot reload (Air)
just templates               # Generate Templ templates (required after template changes)
```

### Building & Running
```bash
just build                   # Full production build (templates + TailwindCSS + Go binary)
just run                     # Run production binary
just run-dev                 # Development mode
just run-prod                # Production mode with optimizations
```

### Testing
```bash
just test                    # Run all tests
go test ./internal/service   # Test specific package
go test -run TestContactForm ./tests/    # Run specific test by name
go test -v ./internal/service/           # Verbose test output
```

### Quality Assurance
```bash
just lint                    # Go vet + golangci-lint (run before commits)
just lint-fix               # Auto-fix linting issues
just arch-lint              # Architectural compliance (DDD/Clean Architecture)
just find-duplicates         # Code duplication analysis
```

### Deployment
```bash
just deploy PROJECT_ID REGION    # Deploy to Google Cloud Run
just deploy-full PROJECT_ID      # Deploy with full quality checks
```

## Architecture Overview

**Clean Architecture + Domain-Driven Design** implementation with clear separation of concerns:

- **`cmd/build/`**: Static site generator for build-time content processing
- **`internal/domain/`**: Core business entities (Experience, Technology, Achievement with quantified metrics)
- **`internal/application/`**: Application services (ContactService for project inquiries)
- **`internal/service/`**: Business logic layer with comprehensive validation
- **`internal/repository/`**: Data access layer with interfaces and implementations
- **`internal/infrastructure/`**: External service implementations (SMTP email service)
- **`internal/handler/`**: HTTP handlers for web endpoints
- **`internal/container/`**: Dependency injection setup (using samber/do)
- **`templates/`**: Templ template files for type-safe HTML generation
- **`tests/`**: End-to-end integration tests with performance thresholds

### Key Technical Patterns

**Dual Container Architecture**: Currently using both portfolio and contact DI containers during architectural transition.

**Quantified Professional Experience**: Domain model includes rich business metrics (test coverage improvements, deployment time reduction, cost savings percentages).

**Type-Safe Templates**: Templ provides compile-time template verification with Go integration.

## Technology Stack

- **Go 1.24.4** with Gin web framework
- **Templ** for HTML templates
- **HTMX** for frontend interactivity  
- **TailwindCSS** with forms and typography plugins
- **Bun** for frontend package management
- **samber/do** for dependency injection
- **Air** for development hot reload
- **golangci-lint** with enterprise-grade configuration (50+ linters)

## Testing Strategy

**Unit Tests**: Service layer with mock repositories (`internal/service/*_test.go`)

**Integration Tests**: Full workflow testing including HTTP handlers and email service (`tests/e2e_integration_test.go`)

**Performance Monitoring**: Response time validation with configurable thresholds

**Architecture Compliance**: go-arch-lint enforces Clean Architecture boundaries

## Development Workflow

1. **Start Development**: `just setup && just templates && just dev`
2. **Make Changes**: Templates auto-regenerate, CSS rebuilds, Go hot-reloads
3. **Test Changes**: `just test` for all tests, `go test -run TestName` for specific tests
4. **Quality Check**: `just lint` before commits
5. **Deploy**: `just deploy-full PROJECT_ID` for production

## Configuration

Environment variables for different deployment modes:
- `ENVIRONMENT`: development/staging/production
- `SERVER_HOST`/`SERVER_PORT`: Server configuration
- `DB_*`: Database configuration  
- Email SMTP settings for contact form functionality

## Quality Standards

**Enterprise-Grade Linting**: Comprehensive golangci-lint with security scanning (gosec), architectural enforcement, and code duplication detection.

**Security**: Path validation, secure SMTP configuration, input sanitization, and proper HTTP timeouts.

**Performance**: Cognitive complexity limits, response time monitoring, and optimized Docker multi-stage builds.
