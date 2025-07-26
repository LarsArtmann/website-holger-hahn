# Holger M. Hahn - Professional Website

A professional website for Holger M. Hahn, Solution Architect & Test Engineer specializing in Digital Assets, Tokenization, and Web3 technologies.

## Features

- **Professional Design**: Clean, modern design targeting financial institutions and enterprise clients
- **Responsive Layout**: Mobile-first responsive design using modern CSS
- **Fast Performance**: Built with Go, Templ, and HTMX for optimal performance
- **Enterprise Focus**: Content optimized for attracting freelance projects from financial institutions

## Technology Stack

- **Backend**: Go with Gin web framework and unified DI container (samber/do)
- **Templates**: Templ for type-safe HTML templates
- **Styling**: TailwindCSS with custom modern responsive design
- **Frontend**: HTMX for interactive features without heavy JavaScript
- **Icons**: Heroicons SVG icons
- **Architecture**: Domain-Driven Design (DDD) with CQRS patterns

## Project Structure

```
.
├── main.go              # Unified web server (portfolio + contact)
├── cmd/build/           # Static site generator
├── internal/            # Application code (DDD structure)
│   ├── application/     # Application services
│   ├── domain/          # Domain entities and interfaces
│   ├── infrastructure/  # External service implementations
│   ├── container/       # Unified DI container
│   ├── handler/         # HTTP handlers
│   ├── service/         # Business logic services
│   └── repository/      # Data access patterns
├── templates/           # Templ template files
│   ├── index.templ      # Main page template
│   └── components.templ # Reusable components
├── static/              # Static assets
│   ├── css/
│   │   ├── input.css    # TailwindCSS source
│   │   └── styles.css   # Compiled CSS
│   └── images/          # Images and icons
├── tests/               # Integration and e2e tests
├── justfile             # Development task automation
└── README.md           # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Templ CLI tool
- Bun (for TailwindCSS builds)
- Just command runner (optional but recommended)

### Installation

1. Install required tools:
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   npm install -g bun
   ```

2. Setup the project:
   ```bash
   just setup  # or: bun install && go mod tidy
   ```

3. Build and run:
   ```bash
   just dev    # Development mode with hot reload
   # or manually:
   just build && go run .
   ```

4. Open your browser and visit: http://localhost:8081

### Available Commands

```bash
just build      # Build CSS and Go binary
just dev        # Development mode with hot reload
just test       # Run all tests
just lint       # Run comprehensive linting
just run        # Run production build
just health     # Check server health
```

## Development

### Template Development

- Templates are written in Templ syntax in the `templates/` directory
- After making changes to `.templ` files, run `templ generate` to compile them
- The server will automatically serve the updated templates

### Styling

- **TailwindCSS**: Source file `static/css/input.css`, compiled to `static/css/styles.css`
- **Build CSS**: Run `bun run build-css` or `just build` to compile TailwindCSS
- **Watch Mode**: Use `bun run watch-css` for development
- **Modern Design**: Utility-first approach with custom components
- **Responsive**: Mobile-first with sm (640px), md (768px), lg (1024px) breakpoints

### Application Features

**Unified Architecture**:
- **Portfolio Display**: Dynamic content from service layer with in-memory repositories
- **Contact Form**: Full contact submission with email notifications
- **Health Monitoring**: Built-in health check endpoints
- **API Endpoints**: RESTful API for technologies, experiences, and services

**Key Sections**:
- Contact information and form submission
- Professional experience with real client work
- Services and expertise areas
- Core technologies and certifications

## Production Deployment

### Docker Deployment (Recommended)

```bash
# Build and deploy with Docker
docker build -t holger-hahn-website .
docker run -p 8080:8080 -e GIN_MODE=release holger-hahn-website
```

### Manual Deployment

```bash
# Production build
just prod-build  # or: bun run build && go build -o holger-hahn-website .

# Run in production mode
ENVIRONMENT=production GIN_MODE=release ./holger-hahn-website
```

### Cloud Run Deployment

```bash
# Deploy to Google Cloud Run
just deploy PROJECT_ID="your-project-id" REGION="us-central1"
```

## Target Audience

This website is designed to attract freelance projects from:
- Financial institutions
- Digital asset companies
- Banks and fintech startups
- Enterprise blockchain projects
- Regulatory compliance projects

## Key Features for Client Acquisition

- **Regulatory Compliance Expertise**: Emphasizes experience with regulated environments
- **Enterprise Client Portfolio**: Showcases work with major financial institutions
- **Security Focus**: Highlights secure, compliant solutions
- **Risk Mitigation**: Demonstrates ability to reduce operational risk
- **Professional Presentation**: Clean, trustworthy design appropriate for financial sector

## License

Private project - All rights reserved.
