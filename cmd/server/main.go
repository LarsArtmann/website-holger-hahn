// Package main provides the HTTP server entry point for Holger Hahn's portfolio website.
// It initializes the dependency injection container, configures the web server,
// and sets up all routes for the portfolio application.
package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"holger-hahn-website/internal/config"
	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/container"
	"holger-hahn-website/internal/handler"
	"holger-hahn-website/internal/service"
	"holger-hahn-website/templates"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DI container.
	di := container.New()
	defer func() {
		if err := di.Shutdown(); err != nil {
			log.Printf("Error shutting down DI container: %v", err)
		}
	}()

	// Get configuration.
	cfg := container.MustGet[*config.Config](di)

	// Set Gin mode based on environment.
	if cfg.Server.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router.
	r := gin.Default()

	// Serve static files.
	r.Static("/static", "./static")

	// Get services from DI container.
	technologyService := container.MustGet[*service.TechnologyService](di)
	experienceService := container.MustGet[*service.ExperienceService](di)
	portfolioService := container.MustGet[*service.PortfolioService](di)

	// Initialize handlers.
	portfolioHandlers := handler.NewPortfolioHandlers(
		technologyService,
		experienceService,
		portfolioService,
	)

	// Setup routes.
	setupRoutes(r, portfolioHandlers)

	// Create HTTP server with configured timeouts.
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	log.Printf("Starting DDD architecture server on %s in %s mode", cfg.Server.Address(), cfg.Server.Environment)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Failed to start server: %v", err)
		// Ensure deferred cleanup runs before exit
		if shutdownErr := di.Shutdown(); shutdownErr != nil {
			log.Printf("Error shutting down DI container: %v", shutdownErr)
		}
		os.Exit(1)
	}
}

// setupRoutes configures all application routes.
func setupRoutes(r *gin.Engine, handlers *handler.PortfolioHandlers) {
	// Health check.
	r.GET("/health", handlers.HealthHandler)

	// Main portfolio page - using existing template for now.
	r.GET("/", func(c *gin.Context) {
		component := templates.Index()

		c.Header("Content-Type", "text/html")

		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			c.JSON(constants.HTTPInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	})

	// API routes for dynamic data.
	api := r.Group("/api/v1")
	{
		api.GET("/technologies", handlers.TechnologiesHandler)
		api.GET("/experiences", handlers.ExperiencesHandler)
		api.GET("/services", handlers.ServicesHandler)
	}
}
