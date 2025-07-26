// Package main provides the unified HTTP server entry point for Holger Hahn's website.
// It combines both portfolio display and contact form functionalities into a single application.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"holger-hahn-website/internal/application"
	"holger-hahn-website/internal/config"
	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/container"
	"holger-hahn-website/internal/handler"
	"holger-hahn-website/internal/infrastructure"
	"holger-hahn-website/internal/service"
	"holger-hahn-website/templates"
)

// ContactHandler handles HTTP requests for contact operations.
type ContactHandler struct {
	contactService *application.ContactService
}

// NewContactHandler creates a new contact handler.
func NewContactHandler(contactService *application.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// SubmitContactForm handles POST /contact requests.
func (h *ContactHandler) SubmitContactForm(c *gin.Context) {
	var req application.ContactFormRequest

	// Bind and validate JSON/form data
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Please check all required fields and try again.",
			"error":   err.Error(),
		})

		return
	}

	// Use application service to handle the request
	ctx := context.Background()

	response, err := h.contactService.SubmitContactForm(ctx, req)
	if err != nil {
		log.Printf("Contact service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Sorry, there was an error processing your request. Please try again later.",
		})

		return
	}

	log.Printf("Contact form submitted successfully - ID: %s", response.ID)
	c.JSON(http.StatusOK, response)
}

// setupRoutes configures all application routes for both portfolio and contact functionality.
func setupRoutes(r *gin.Engine, portfolioHandlers *handler.PortfolioHandlers, contactHandler *ContactHandler) {
	// Serve static files
	r.Static("/static", "./static")

	// Main portfolio page
	r.GET("/", func(c *gin.Context) {
		component := templates.Index()

		c.Header("Content-Type", "text/html")

		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			c.JSON(constants.HTTPInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	})

	// Contact form API endpoint
	r.POST("/contact", contactHandler.SubmitContactForm)

	// Health check endpoint
	r.GET("/health", portfolioHandlers.HealthHandler)

	// Portfolio API routes for dynamic data
	api := r.Group("/api/v1")
	{
		api.GET("/technologies", portfolioHandlers.TechnologiesHandler)
		api.GET("/experiences", portfolioHandlers.ExperiencesHandler)
		api.GET("/services", portfolioHandlers.ServicesHandler)
	}
}

func main() {
	// Initialize unified DI container (using portfolio app's container system)
	di := container.New()

	// Get configuration
	cfg := container.MustGet[*config.Config](di)

	// Set Gin mode based on environment
	if cfg.Server.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Get portfolio services from DI container
	technologyService := container.MustGet[*service.TechnologyService](di)
	experienceService := container.MustGet[*service.ExperienceService](di)
	portfolioService := container.MustGet[*service.PortfolioService](di)

	// Initialize portfolio handlers
	portfolioHandlers := handler.NewPortfolioHandlers(
		technologyService,
		experienceService,
		portfolioService,
	)

	// Setup contact service using infrastructure container (temporary dual container approach)
	contactInjector := infrastructure.SetupContainer()
	contactService := do.MustInvoke[*application.ContactService](contactInjector)
	contactHandler := NewContactHandler(contactService)

	// Setup all routes (portfolio + contact)
	setupRoutes(r, portfolioHandlers, contactHandler)

	// Create HTTP server with configured timeouts
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	log.Printf("🚀 Unified Holger Hahn website server starting on %s in %s mode", cfg.Server.Address(), cfg.Server.Environment)
	log.Println("📧 Contact form endpoint: POST /contact")
	log.Println("🏥 Health check: GET /health")
	log.Println("🔧 Portfolio API: GET /api/v1/technologies, /api/v1/experiences, /api/v1/services")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Failed to start server: %v", err)
		// Cleanup DI containers before exit
		if shutdownErr := di.Shutdown(); shutdownErr != nil {
			log.Printf("Error shutting down DI container: %v", shutdownErr)
		}

		os.Exit(constants.ExitFailure)
	}

	// Normal shutdown - cleanup DI containers
	if err := di.Shutdown(); err != nil {
		log.Printf("Error shutting down DI container: %v", err)
	}
}
