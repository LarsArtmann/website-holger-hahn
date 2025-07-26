# Report about missing/under-specified/confusing information

Date: 2025-07-26T06:43:02+02:00

I was asked to perform:
Comprehensive code quality improvement using Pareto principle, then brutal honesty assessment followed by comprehensive architectural improvement plan with GitHub issue management and documentation.

I was given these context information's:
- List of required libraries: gin-gonic/gin, spf13/viper, a-h/templ, htmx, fe3dback/go-arch-lint, samber/lo, samber/mo, samber/do, sqlc-dev/sqlc, onsi/ginkgo, charmbracelet/fang, OpenTelemetry, LarsArtmann/uniflow
- Architecture patterns: Event-Sourcing, DDD, CQRS, Composition over inheritance, FP patterns, Layered Architecture, EDA, Railway Oriented Programming, BDD, TDD
- Requirement to not reinvent the wheel and leverage existing libraries
- Instructions to identify and integrate ghost systems
- Target legacy code reduction to 0

I was missing these information:
- Clear definition of what constitutes a "ghost system" vs legitimate architectural separation
- Priority ranking between different architectural patterns when they conflict
- Specific business requirements driving the need for event sourcing vs simpler patterns
- Clear integration strategy between existing contact form system and portfolio system
- Database schema understanding and migration strategy from in-memory to persistent storage
- Performance requirements that would justify the complexity of full DDD/CQRS/Event Sourcing

I was confused by:
- Existence of two separate main.go applications (main.go and cmd/server/main.go) - unclear if this was intentional or accidental
- Disconnect between specified required libraries and actual go.mod dependencies
- Extensive sqlc generated code that appears unused - unclear if this is legacy or intended
- Mixed architectural patterns (some proper DDD in domain layer, some simple structs elsewhere)
- Unclear customer value proposition for implementing full event sourcing for what appears to be a relatively simple portfolio website

What I wish for the future is:
- Clear architectural decision records (ADRs) documenting why specific patterns were chosen
- Explicit business requirements that drive architectural complexity decisions
- Clear integration strategy documentation when multiple systems exist
- Library usage guidelines that specify when to use which patterns
- Clear definition of "production ready" vs "over-engineered" for different contexts
- Step-by-step migration guides from simpler to more complex architectures
- Clear customer value metrics for each architectural decision

Best regards,
Claude (Sonnet 4)