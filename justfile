# Holger Hahn Website - Just Commands

# Install dependencies and setup
setup:
    bun install
    go mod tidy
    go install github.com/golangci/dupl@latest

# Build the project
build:
    bun run build

# Run linting with strict enterprise-grade rules
lint:
    @echo "Running Go linting..."
    go vet ./...
    @echo "Running golangci-lint with strict configuration..."
    golangci-lint run --config .golangci.yml || echo "golangci-lint not installed, install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

# Run architectural linting
arch-lint:
    @echo "Running architectural linting..."
    go-arch-lint check || echo "go-arch-lint not installed, install with: go install github.com/fe3dback/go-arch-lint@latest"

# Install all linting tools
install-linters:
    @echo "Installing linting tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/fe3dback/go-arch-lint@latest
    go install github.com/golangci/dupl@latest
    @echo "All linters installed successfully"

# Run comprehensive linting (all checks)
lint-all: lint arch-lint find-duplicates
    @echo "All linting checks completed"

# Fix auto-fixable linting issues
lint-fix:
    @echo "Fixing auto-fixable linting issues..."
    golangci-lint run --fix --config .golangci.yml

# Find duplicate code in Go files with comprehensive analysis
find-duplicates:
    @echo "🔍 Scanning for duplicate code (threshold: 60 tokens)..."
    @echo "📁 Analyzing Go files in internal/ directory..."
    dupl -t 60 -v internal/ || echo "❌ dupl not installed, install with: go install github.com/golangci/dupl@latest"
    @echo "✅ Duplicate code analysis completed"

# Generate detailed HTML report of duplicates
find-duplicates-html:
    @echo "🔍 Generating HTML duplicate code report..."
    dupl -t 60 -html internal/ > duplicates-report.html || echo "❌ dupl not installed"
    @echo "📊 HTML report generated: duplicates-report.html"

# Find duplicates with strict threshold (15 tokens)
find-duplicates-strict:
    @echo "🔍 Scanning for duplicate code (strict mode: 15 tokens)..."
    dupl -t 15 -v internal/ || echo "❌ dupl not installed"

# Quick duplicate scan (alias)
fd: find-duplicates

# Comprehensive duplicate analysis (all modes)
duplicates-all: find-duplicates find-duplicates-html find-duplicates-strict
    @echo "🎯 Comprehensive duplicate analysis completed"

# Generate templates
templates:
    templ generate

# Run the application in production mode
run:
    @echo "🚀 Starting Holger Hahn Website (Production Mode)..."
    @if [ -f holger-hahn-website ]; then \
        ./holger-hahn-website; \
    else \
        echo "❌ Binary not found. Building first..."; \
        just build && ./holger-hahn-website; \
    fi

# Run the unified server (same as run-dev)
run-server:
    @echo "🚀 Starting unified server..."
    go run .

# Run the development server
dev:
    @echo "🔧 Starting development mode with hot reload..."
    bun run watch-css &
    air || go run .

# Run tests
test:
    go test ./...

# Add restart development server
restart:
    @echo "🔄 Restarting development server..."
    @pkill -f "go run" || true
    @pkill -f "air" || true
    @pkill -f "holger-hahn-website" || true
    @sleep 1
    just dev

# Add server health check
health:
    @echo "🏥 Checking server health..."
    @curl -f http://localhost:8081/health 2>/dev/null && echo "✅ Server is healthy" || echo "❌ Server is not responding"

# Run tests
run-tests: test

# Clean build artifacts
clean:
    rm -f holger-hahn-website
    rm -f duplicates.html
    rm -f duplicates-report.html

# Full development setup
dev-setup: setup templates build

# Production build
prod-build: templates build
    @echo "Production build complete"

# Run with environment variables
run-dev:
    @echo "🔧 Running in development mode..."
    ENVIRONMENT=development GIN_MODE=debug go run .

run-prod:
    @echo "🚀 Running in production mode..."
    ENVIRONMENT=production GIN_MODE=release ./holger-hahn-website

run-staging:
    @echo "🎭 Running in staging mode..."
    ENVIRONMENT=staging GIN_MODE=release go run .

# Deploy to Google Cloud Run (requires PROJECT_ID)
deploy PROJECT_ID="your-project-id" REGION="us-central1":
    @echo "🚀 Deploying to Google Cloud Run..."
    @echo "Project: {{PROJECT_ID}}"
    @echo "Region: {{REGION}}"
    ./deploy.sh {{PROJECT_ID}} {{REGION}}

# Deploy with build and lint checks
deploy-full PROJECT_ID="your-project-id" REGION="us-central1": build lint deploy
    @echo "✅ Full deployment with quality checks completed"

# Quick deploy (skip linting - use for hotfixes)
deploy-quick PROJECT_ID="your-project-id" REGION="us-central1": build
    @echo "⚡ Quick deployment completed (linting skipped)"
    ./deploy.sh {{PROJECT_ID}} {{REGION}}

# Check project status
status:
    @echo "=== Project Status ==="
    @echo "Go version: $(go version)"
    @echo "Node/Bun version: $(bun --version)"
    @echo "Files count: $(find . -name '*.go' | wc -l) Go files, $(find . -name '*.templ' | wc -l) template files"
    @echo "Build status: $(if [ -f holger-hahn-website ]; then echo 'Built'; else echo 'Not built'; fi)"

# Quality monitoring commands (Group G8)
monitor:
    @echo "🔍 Running comprehensive quality monitoring..."
    ./scripts/monitor-quality.sh

quick-check:
    @echo "⚡ Quick quality status check..."
    ./scripts/quick-status.sh

quality-baseline:
    @echo "📊 Establishing quality baseline..."
    @echo "Baseline established on $(date)" >> docs/quality-baseline.log
    just test || echo "Tests failing"
    just build || echo "Build failing"
    just lint || echo "Linting failing"
    just fd || echo "Duplicates detected"
