# Dockerfile for Holger Hahn Website - Cloud Run Deployment
# Multi-stage build for optimal container size and security

# Stage 1: Build stage
FROM node:18-alpine AS node-builder

WORKDIR /app

# Copy package files for Node.js dependencies
COPY package.json bun.lockb* ./

# Install bun for faster builds
RUN npm install -g bun

# Install dependencies
RUN bun install

# Copy source files needed for CSS build
COPY . .

# Build CSS with TailwindCSS
RUN bun run build-css

# Stage 2: Go build stage
FROM golang:1.21-alpine AS go-builder

# Install templ for template generation
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built CSS from node-builder stage
COPY --from=node-builder /app/static/css/styles.css ./static/css/styles.css

# Generate templates
RUN templ generate

# Build Go application with optimizations for Cloud Run
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -tags netgo \
    -o holger-hahn-website \
    .

# Stage 3: Final runtime stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Copy binary from builder stage
COPY --from=go-builder /app/holger-hahn-website ./

# Copy static assets
COPY --from=go-builder /app/static ./static

# Create necessary directories and set permissions
RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (Cloud Run will set PORT environment variable)
EXPOSE 8080

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-8080}/health || exit 1

# Environment variables for production
ENV GIN_MODE=release
ENV ENVIRONMENT=production
ENV PORT=8080

# Run the application
CMD ["./holger-hahn-website"]