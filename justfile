# Holger Hahn Website - Just Commands

# Install dependencies and setup
setup:
    bun install
    go mod tidy
    go install github.com/golangci/dupl@latest

# Build the project
build:
    bun run build

# Run linting
lint:
    @echo "Running Go linting..."
    go vet ./...
    @echo "Running golangci-lint..."
    golangci-lint run || echo "golangci-lint not installed, install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

# Find duplicate code in Go files
find-duplicates:
    @echo "Scanning for duplicate code..."
    dupl -t 15 || echo "dupl not installed, install with: go install github.com/golangci/dupl@latest"

# Alias for find-duplicates
fd: find-duplicates

# Generate templates
templates:
    templ generate

# Run the development server
dev:
    bun run watch-css &
    air || go run .

# Run tests
test:
    go test ./...

# Clean build artifacts
clean:
    rm -f holger-hahn-website
    rm -f duplicates.html

# Full development setup
dev-setup: setup templates build

# Production build
prod-build: templates build
    @echo "Production build complete"

# Check project status
status:
    @echo "=== Project Status ==="
    @echo "Go version: $(go version)"
    @echo "Node/Bun version: $(bun --version)"
    @echo "Files count: $(find . -name '*.go' | wc -l) Go files, $(find . -name '*.templ' | wc -l) template files"
    @echo "Build status: $(if [ -f holger-hahn-website ]; then echo 'Built'; else echo 'Not built'; fi)"