package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"holger-hahn-website/internal/config"
	"holger-hahn-website/internal/container"
	"holger-hahn-website/internal/handler"
	"holger-hahn-website/internal/service"
	"holger-hahn-website/templates"
)

func main() {
	// Initialize DI container
	di := container.New()
	defer di.Shutdown()

	// Get configuration
	cfg := container.MustGet[*config.Config](di)

	// Set Gin mode based on environment
	if cfg.Server.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Get services from DI container
	technologyService := container.MustGet[*service.TechnologyService](di)
	experienceService := container.MustGet[*service.ExperienceService](di)
	portfolioService := container.MustGet[*service.PortfolioService](di)

	// Initialize handlers
	portfolioHandlers := handler.NewPortfolioHandlers(
		technologyService,
		experienceService,
		portfolioService,
	)

	// Setup routes
	setupRoutes(r, portfolioHandlers)

	// Create HTTP server with configured timeouts
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	log.Printf("Starting DDD architecture server on %s in %s mode", cfg.Server.Address(), cfg.Server.Environment)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes configures all application routes
func setupRoutes(r *gin.Engine, handlers *handler.PortfolioHandlers) {
	// Health check
	r.GET("/health", handlers.HealthHandler)

	// Main portfolio page - using existing template for now
	r.GET("/", func(c *gin.Context) {
		component := templates.Index()
		c.Header("Content-Type", "text/html")
		component.Render(c.Request.Context(), c.Writer)
	})

	// API routes for dynamic data
	api := r.Group("/api/v1")
	{
		api.GET("/technologies", handlers.TechnologiesHandler)
		api.GET("/experiences", handlers.ExperiencesHandler)
		api.GET("/services", handlers.ServicesHandler)
	}
}