# Report about missing/under-specified/confusing information

Date: 2025-07-26T07:05:00+02:00

I was asked to perform:
Comprehensive code quality improvement with brutal honesty assessment, create detailed execution plans (30-100min tasks, then 12min micro-tasks), leverage existing libraries instead of reinventing, identify and integrate ghost systems, reduce legacy code to 0, use incremental git commits, manage GitHub issues, create extensive documentation, and focus on customer value delivery.

I was given these context information's:
- List of preferred libraries: gin-gonic/gin, spf13/viper, a-h/templ, htmx, samber/lo, samber/mo, samber/do, sqlc-dev/sqlc, onsi/ginkgo, charmbracelet/fang, OpenTelemetry, LarsArtmann/uniflow
- Architecture patterns: Event-Sourcing, DDD, CQRS, Composition over inheritance, FP patterns, Layered Architecture, EDA, Railway Oriented Programming, BDD, TDD
- Instructions to use existing code and not reinvent the wheel
- Requirement to identify ghost systems and integrate them
- Target of 0 legacy code
- Emphasis on customer value creation
- Request for brutal honesty about failures and stupid decisions

I was missing these information:
- **Business requirements clarity** - What specific customer problems does the portfolio website need to solve?
- **Performance requirements** - What response times and throughput are actually needed?
- **Budget/time constraints** - Is implementing full Event Sourcing/CQRS justified for a portfolio website?
- **User expectations** - How complex should the contact form and portfolio display actually be?
- **Production requirements** - What's the expected traffic and data volume?
- **Integration requirements** - Does this need to integrate with other systems?
- **Compliance requirements** - What specific financial regulations apply to Holger's business?

I was confused by:
- **Scope mismatch** - The architectural complexity (Event Sourcing, CQRS, DDD) seems over-engineered for a personal portfolio website
- **Priority contradiction** - Focus on customer value vs implementing complex architectural patterns that may not deliver customer value
- **Time allocation expectation** - Whether to spend time on execution vs documentation when both were requested
- **Definition of "ghost systems"** - Whether unused generated code counts as ghost systems vs intentional future preparation
- **Legacy code definition** - What constitutes legacy code vs working code that could be improved
- **Customer value measurement** - How to quantify customer value for architectural patterns vs functional features

What I wish for the future is:
- **Clear business case analysis** before implementing complex patterns - justify Event Sourcing/CQRS with specific business needs
- **Customer value metrics definition** - How to measure success for architectural decisions
- **Explicit scope boundaries** - What's the minimum viable architecture vs gold-plating
- **Time allocation guidance** - When to prioritize execution vs documentation vs analysis
- **Ghost system definition** - Clear criteria for what constitutes a ghost system vs future preparation
- **ROI analysis framework** - Cost/benefit analysis for architectural complexity decisions
- **Progressive complexity approach** - Start simple, add complexity only when business requirements justify it
- **Customer interview data** - What do actual users of Holger's website need?
- **Performance baseline requirements** - Specific targets rather than generic "make it fast"
- **Integration roadmap** - What future integrations are planned that would justify current complexity?

Best regards,
Claude (Sonnet 4)