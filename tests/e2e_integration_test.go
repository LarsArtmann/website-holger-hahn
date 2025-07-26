// Package tests provides end-to-end integration tests for the portfolio website.
// It contains comprehensive tests that verify the complete functionality of the application
// including HTTP handlers, business logic, and infrastructure integration.
package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"holger-hahn-website/internal/application"
	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/container"
	"holger-hahn-website/templates"
)

// TestEndToEndContactFlow tests the complete contact form workflow.
func TestEndToEndContactFlow(t *testing.T) {
	// Setup unified dependency injection container
	di := container.New()
	defer func() {
		if err := di.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown DI container: %v", err)
		}
	}()

	contactService := container.MustGet[*application.ContactService](di)

	// Create contact handler
	contactHandler := &ContactHandler{contactService: contactService}

	// Setup Gin router in test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Serve static files
	router.Static("/static", "../static")

	// Routes
	router.GET("/", func(c *gin.Context) {
		component := templates.Index()

		c.Header("Content-Type", "text/html")

		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			c.JSON(constants.HTTPInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	})

	router.POST("/contact", contactHandler.SubmitContactForm)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(constants.HTTPOKStatus, gin.H{"status": "healthy"})
	})

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", constants.HTTPOKStatus, resp.StatusCode)
		}

		var healthResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
			t.Errorf("Failed to decode health response: %v", err)
		}

		if healthResp["status"] != "healthy" {
			t.Errorf("Expected status 'healthy', got '%v'", healthResp["status"])
		}
	})

	t.Run("Homepage Loads", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("Homepage request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", constants.HTTPOKStatus, resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Expected HTML content type, got %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		bodyStr := string(body)

		// Verify essential elements are present
		essentialElements := []string{
			"Holger M. Hahn",
			"Digital Assets",
			"Blockchain Solutions",
			"contact form",
			"/static/css/styles.css",
		}

		for _, element := range essentialElements {
			if !strings.Contains(bodyStr, element) {
				t.Errorf("Homepage missing essential element: %s", element)
			}
		}
	})

	t.Run("Contact Form Submission Success", func(t *testing.T) {
		contactData := map[string]interface{}{
			"name":    "John Doe",
			"company": "Test Corp",
			"email":   "john.doe@example.com",
			"project": "I need help with implementing a blockchain solution for our financial services company. We're looking for expertise in regulatory compliance and digital asset custody.",
		}

		jsonData, err := json.Marshal(contactData)
		if err != nil {
			t.Fatalf("Failed to marshal contact data: %v", err)
		}

		resp, err := http.Post(server.URL+"/contact", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Contact form submission failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body))
		}

		var response application.ContactFormResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("Expected success=true, got %v", response.Success)
		}

		if response.ID == "" {
			t.Errorf("Expected non-empty ID in response")
		}

		if !strings.Contains(response.Message, "Thank you") {
			t.Errorf("Expected thank you message, got: %s", response.Message)
		}
	})

	t.Run("Contact Form Validation Errors", func(t *testing.T) {
		invalidTestCases := []struct {
			data map[string]interface{}
			name string
		}{
			{
				name: "Missing Name",
				data: map[string]interface{}{
					"company": "Test Corp",
					"email":   "john@example.com",
					"project": "This is a test project description that is long enough.",
				},
			},
			{
				name: "Invalid Email",
				data: map[string]interface{}{
					"name":    "John Doe",
					"company": "Test Corp",
					"email":   "invalid-email",
					"project": "This is a test project description that is long enough.",
				},
			},
			{
				name: "Project Too Short",
				data: map[string]interface{}{
					"name":    "John Doe",
					"company": "Test Corp",
					"email":   "john@example.com",
					"project": "Short",
				},
			},
		}

		for _, tc := range invalidTestCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonData, err := json.Marshal(tc.data)
				if err != nil {
					t.Fatalf("Failed to marshal test data: %v", err)
				}

				resp, err := http.Post(server.URL+"/contact", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					t.Fatalf("Request failed: %v", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status %d, got %d", constants.HTTPBadRequestStatus, resp.StatusCode)
				}

				var errorResp map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}

				if errorResp["success"] != false {
					t.Errorf("Expected success=false for validation error")
				}
			})
		}
	})
}

// ContactHandler for testing.
type ContactHandler struct {
	contactService *application.ContactService
}

func (h *ContactHandler) SubmitContactForm(c *gin.Context) {
	var req application.ContactFormRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Please check all required fields and try again.",
			"error":   err.Error(),
		})

		return
	}

	ctx := context.Background()

	response, err := h.contactService.SubmitContactForm(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Sorry, there was an error processing your request. Please try again later.",
		})

		return
	}

	c.JSON(http.StatusOK, response)
}

// TestSystemIntegration tests all system components working together.
func TestSystemIntegration(t *testing.T) {
	t.Run("Template Generation", func(t *testing.T) {
		// Test that templates can be generated without errors
		component := templates.Index()
		if component == nil {
			t.Error("Failed to create Index template component")
		}

		// Test rendering
		var buf bytes.Buffer

		ctx := t.Context()
		if err := component.Render(ctx, &buf); err != nil {
			t.Errorf("Failed to render Index template: %v", err)
		}

		rendered := buf.String()
		if len(rendered) == 0 {
			t.Error("Template rendered to empty string")
		}

		// Verify HTML structure
		if !strings.Contains(strings.ToLower(rendered), "<!doctype html>") {
			t.Error("Template missing DOCTYPE declaration")
		}

		if !strings.Contains(rendered, "</html>") {
			t.Error("Template missing closing html tag")
		}
	})

	t.Run("Dependency Injection", func(t *testing.T) {
		// Test that unified DI container can be set up
		di := container.New()
		if di == nil {
			t.Fatal("Failed to setup unified DI container")
		}

		// Test that services can be resolved
		contactService := container.MustGet[*application.ContactService](di)
		if contactService == nil {
			t.Error("Failed to resolve ContactService from unified DI container")
		}

		// Cleanup container
		if err := di.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown DI container: %v", err)
		}
	})

	t.Run("Static File Serving", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		router.Static("/static", "../static")

		server := httptest.NewServer(router)
		defer server.Close()

		// Test CSS file serving
		resp, err := http.Get(server.URL + "/static/css/styles.css")
		if err != nil {
			t.Fatalf("Failed to get CSS file: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for CSS file, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/css") {
			t.Errorf("Expected CSS content type, got %s", contentType)
		}
	})
}

// TestPerformanceBasic tests basic performance expectations.
func TestPerformanceBasic(t *testing.T) {
	di := container.New()
	defer func() {
		if err := di.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown DI container: %v", err)
		}
	}()

	contactService := container.MustGet[*application.ContactService](di)
	contactHandler := &ContactHandler{contactService: contactService}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Static("/static", "../static")
	router.GET("/", func(c *gin.Context) {
		component := templates.Index()

		c.Header("Content-Type", "text/html")

		if err := component.Render(c.Request.Context(), c.Writer); err != nil {
			c.JSON(constants.HTTPInternalServerError, gin.H{"error": "Failed to render template"})
			return
		}
	})
	router.POST("/contact", contactHandler.SubmitContactForm)

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Homepage Response Time", func(t *testing.T) {
		start := time.Now()
		resp, err := http.Get(server.URL + "/")
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Homepage request failed: %v", err)
		}

		defer resp.Body.Close()

		// Homepage should load in under 100ms for local testing
		if elapsed > constants.FastResponseThreshold {
			t.Logf("Homepage took %v to load (acceptable but monitoring)", elapsed)
		}

		if elapsed > constants.SlowResponseThreshold {
			t.Errorf("Homepage took too long to load: %v", elapsed)
		}
	})

	t.Run("Contact Form Response Time", func(t *testing.T) {
		contactData := map[string]interface{}{
			"name":    "Performance Test",
			"company": "Test Corp",
			"email":   "perf@example.com",
			"project": "This is a performance test submission to measure response time.",
		}

		jsonData, _ := json.Marshal(contactData)

		start := time.Now()
		resp, err := http.Post(server.URL+"/contact", "application/json", bytes.NewBuffer(jsonData))
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Contact form submission failed: %v", err)
		}

		defer resp.Body.Close()

		// Contact form should process in under 200ms for local testing
		if elapsed > constants.AcceptableResponseThreshold {
			t.Logf("Contact form took %v to process (acceptable but monitoring)", elapsed)
		}

		if elapsed > constants.UnacceptableResponseThreshold {
			t.Errorf("Contact form took too long to process: %v", elapsed)
		}
	})
}
