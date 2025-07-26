package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"holger-hahn-website/internal/application"
	"holger-hahn-website/internal/infrastructure"
	"holger-hahn-website/templates"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
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

	// Return success response
	c.JSON(http.StatusOK, response)
}

func main() {
	// Setup dependency injection container
	injector := infrastructure.SetupContainer()

	// Get contact service from container
	contactService := do.MustInvoke[*application.ContactService](injector)

	// Create handlers
	contactHandler := NewContactHandler(contactService)

	// Setup Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Routes
	r.GET("/", func(c *gin.Context) {
		component := templates.Index()

		c.Header("Content-Type", "text/html")

		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			c.JSON(500, gin.H{"error": "Failed to render template"})
			return
		}
	})

	// Contact form API endpoint
	r.POST("/contact", contactHandler.SubmitContactForm)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Get port from environment (Cloud Run sets PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default for local development
	}

	log.Printf("Server starting on :%s", port)
	log.Println("Contact form endpoint: POST /contact")
	log.Println("Health check: GET /health")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
