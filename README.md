# Holger M. Hahn - Professional Website

A professional website for Holger M. Hahn, Solution Architect & Test Engineer specializing in Digital Assets, Tokenization, and Web3 technologies.

## Features

- **Professional Design**: Clean, modern design targeting financial institutions and enterprise clients
- **Responsive Layout**: Mobile-first responsive design using modern CSS
- **Fast Performance**: Built with Go, Templ, and HTMX for optimal performance
- **Enterprise Focus**: Content optimized for attracting freelance projects from financial institutions

## Technology Stack

- **Backend**: Go with Gin web framework
- **Templates**: Templ for type-safe HTML templates
- **Styling**: Custom CSS with modern responsive design
- **Frontend**: HTMX for interactive features
- **Icons**: Heroicons SVG icons

## Project Structure

```
.
├── main.go              # Go web server
├── templates/           # Templ template files
│   ├── index.templ      # Main page template
│   └── components.templ # Reusable components
├── static/              # Static assets
│   ├── css/
│   │   ├── input.css    # Source CSS file
│   │   └── styles.css   # Compiled CSS
│   └── js/              # JavaScript files
├── run.sh               # Development server script
└── README.md           # This file
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Templ CLI tool

### Installation

1. Install Templ CLI:
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate templates:
   ```bash
   templ generate
   ```

4. Run the development server:
   ```bash
   ./run.sh
   ```

   Or manually:
   ```bash
   go run main.go
   ```

5. Open your browser and visit: http://localhost:8080

## Development

### Template Development

- Templates are written in Templ syntax in the `templates/` directory
- After making changes to `.templ` files, run `templ generate` to compile them
- The server will automatically serve the updated templates

### Styling

- CSS is located in `static/css/styles.css`
- The design uses a modern utility-first approach
- Responsive breakpoints: mobile-first with sm (640px), md (768px), lg (1024px)

### Content Updates

Key sections to customize:
- Contact information (email, LinkedIn)
- Professional experience details
- Services offered
- Educational background
- Certifications

## Production Deployment

1. Set Gin to release mode:
   ```bash
   export GIN_MODE=release
   ```

2. Build the application:
   ```bash
   go build -o website main.go
   ```

3. Run the compiled binary:
   ```bash
   ./website
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