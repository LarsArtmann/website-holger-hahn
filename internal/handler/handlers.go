// Package handler provides HTTP request handlers for the portfolio website.
// It contains Gin HTTP handlers that process incoming requests, coordinate with business services,
// and render appropriate responses for both API endpoints and web pages.
package handler

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/domain"
	"holger-hahn-website/internal/service"
	"holger-hahn-website/templates"
)

// PortfolioHandlers contains all HTTP handlers for the portfolio application.
type PortfolioHandlers struct {
	technologyService *service.TechnologyService
	experienceService *service.ExperienceService
	portfolioService  *service.PortfolioService
}

// NewPortfolioHandlers creates a new portfolio handlers instance.
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

// HomeHandler renders the main portfolio page with graceful error handling.
func (h *PortfolioHandlers) HomeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Get current experiences with fallback on error
	experiences, err := h.experienceService.GetCurrentExperiences(ctx)
	if err != nil {
		// Log error for debugging but continue with empty data
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Warning: Failed to load experiences: %v\n", err)))
		experiences = []*domain.Experience{} // Fallback to empty slice
	}

	// Get active services with fallback on error
	services, err := h.portfolioService.GetActiveServices(ctx)
	if err != nil {
		// Log error for debugging but continue with empty data
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Warning: Failed to load services: %v\n", err)))
		services = []*domain.Service{} // Fallback to empty slice
	}

	// Get technologies with fallback on error
	expertLevel := "expert"
	techFilter := service.TechnologyFilter{
		OrderBy:  &[]string{"name"}[0],
		OrderDir: &[]string{"asc"}[0],
	}

	technologies, err := h.technologyService.ListTechnologies(ctx, techFilter)
	if err != nil {
		// Log error for debugging but continue with empty data
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Warning: Failed to load technologies: %v\n", err)))
		technologies = []*domain.Technology{} // Fallback to empty slice
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

	// Render template with dynamic data (even if some data failed to load)
	component := templates.IndexWithData(experiences)

	c.Header("Content-Type", "text/html")

	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		// Template rendering failure is a critical error that should return HTTP error
		gin.DefaultWriter.Write([]byte(fmt.Sprintf("Error: Failed to render template: %v\n", err)))
		c.JSON(constants.HTTPInternalServerError, gin.H{"error": domain.ErrRenderTemplate.Error()})
		return
	}
}

// HealthHandler provides a health check endpoint.
func (h *PortfolioHandlers) HealthHandler(c *gin.Context) {
	c.JSON(constants.HTTPOKStatus, gin.H{
		"status":  "healthy",
		"version": constants.AppVersion,
	})
}

// TechnologiesHandler returns JSON list of technologies.
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
		c.JSON(constants.HTTPInternalServerError, gin.H{"error": domain.ErrLoadTechnologies.Error()})
		return
	}

	c.JSON(constants.HTTPOKStatus, technologies)
}

// ExperiencesHandler returns JSON list of experiences.
func (h *PortfolioHandlers) ExperiencesHandler(c *gin.Context) {
	ctx := context.Background()

	experiences, err := h.experienceService.ListExperiences(ctx, service.ExperienceFilter{})
	if err != nil {
		c.JSON(constants.HTTPInternalServerError, gin.H{"error": domain.ErrLoadExperiences.Error()})
		return
	}

	c.JSON(constants.HTTPOKStatus, experiences)
}

// ServicesHandler returns JSON list of services.
func (h *PortfolioHandlers) ServicesHandler(c *gin.Context) {
	ctx := context.Background()

	services, err := h.portfolioService.ListServices(ctx, service.ServiceFilter{})
	if err != nil {
		c.JSON(constants.HTTPInternalServerError, gin.H{"error": domain.ErrLoadServices.Error()})
		return
	}

	c.JSON(constants.HTTPOKStatus, services)
}

// PortfolioData represents the data structure passed to templates.
type PortfolioData struct {
	Experiences  interface{} `json:"experiences"`
	Services     interface{} `json:"services"`
	Technologies []string    `json:"technologies"`
}
