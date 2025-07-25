package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"holger-hahn-website/internal/service"
	"holger-hahn-website/templates"
)

// PortfolioHandlers contains all HTTP handlers for the portfolio application
type PortfolioHandlers struct {
	technologyService *service.TechnologyService
	experienceService *service.ExperienceService
	portfolioService  *service.PortfolioService
}

// NewPortfolioHandlers creates a new portfolio handlers instance
func NewPortfolioHandlers(
	technologyService *service.TechnologyService,
	experienceService *service.ExperienceService,
	portfolioService *service.PortfolioService,
) *PortfolioHandlers {
	return &PortfolioHandlers{
		technologyService: technologyService,
		experienceService: experienceService,
		portfolioService:  portfolioService,
	}
}

// HomeHandler renders the main portfolio page
func (h *PortfolioHandlers) HomeHandler(c *gin.Context) {
	ctx := context.Background()
	
	// Get current experiences
	experiences, err := h.experienceService.GetCurrentExperiences(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load experiences"})
		return
	}
	
	// Get active services
	services, err := h.portfolioService.GetActiveServices(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load services"})
		return
	}
	
	// Get technologies (example: showing expert level technologies)
	expertLevel := "expert"
	techFilter := service.TechnologyFilter{
		OrderBy:  &[]string{"name"}[0],
		OrderDir: &[]string{"asc"}[0],
	}
	technologies, err := h.technologyService.ListTechnologies(ctx, techFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load technologies"})
		return
	}
	
	// Filter for expert technologies for the front page
	var expertTechnologies []string
	for _, tech := range technologies {
		if string(tech.Level) == expertLevel {
			expertTechnologies = append(expertTechnologies, tech.Name)
		}
	}
	
	// Create portfolio data structure (will be used when template is updated)
	_ = PortfolioData{
		Experiences:  experiences,
		Services:     services,
		Technologies: expertTechnologies,
	}
	
	// Render template with data (using existing template for now)
	// TODO: Create IndexWithData template that accepts portfolioData
	component := templates.Index()
	c.Header("Content-Type", "text/html")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.JSON(500, gin.H{"error": "Failed to render template"})
		return
	}
}

// HealthHandler provides a health check endpoint
func (h *PortfolioHandlers) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"version": "1.0.0",
	})
}

// TechnologiesHandler returns JSON list of technologies
func (h *PortfolioHandlers) TechnologiesHandler(c *gin.Context) {
	ctx := context.Background()
	
	category := c.Query("category")
	level := c.Query("level")
	
	filter := service.TechnologyFilter{}
	if category != "" {
		filter.Category = &category
	}
	// TODO: Implement level filtering when level != ""
	// Currently level filtering is not implemented
	_ = level // Acknowledge the parameter to avoid unused variable warning
	
	technologies, err := h.technologyService.ListTechnologies(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load technologies"})
		return
	}
	
	c.JSON(http.StatusOK, technologies)
}

// ExperiencesHandler returns JSON list of experiences
func (h *PortfolioHandlers) ExperiencesHandler(c *gin.Context) {
	ctx := context.Background()
	
	experiences, err := h.experienceService.ListExperiences(ctx, service.ExperienceFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load experiences"})
		return
	}
	
	c.JSON(http.StatusOK, experiences)
}

// ServicesHandler returns JSON list of services
func (h *PortfolioHandlers) ServicesHandler(c *gin.Context) {
	ctx := context.Background()
	
	services, err := h.portfolioService.ListServices(ctx, service.ServiceFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load services"})
		return
	}
	
	c.JSON(http.StatusOK, services)
}

// PortfolioData represents the data structure passed to templates
type PortfolioData struct {
	Experiences  interface{} `json:"experiences"`
	Services     interface{} `json:"services"`
	Technologies []string    `json:"technologies"`
}