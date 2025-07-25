package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"holger-hahn-website/templates"
)

// ContactRequest represents the contact form data
type ContactRequest struct {
	Name    string `json:"name" form:"name" binding:"required,min=2,max=100"`
	Company string `json:"company" form:"company" binding:"max=100"`
	Email   string `json:"email" form:"email" binding:"required,email"`
	Project string `json:"project" form:"project" binding:"required,min=10,max=2000"`
}

// ContactResponse represents the API response
type ContactResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	ToEmail      string
}

func getEmailConfig() *EmailConfig {
	// Default to environment variables with fallbacks for development
	smtpPortStr := getEnvOrDefault("SMTP_PORT", "587")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		smtpPort = 587
	}

	return &EmailConfig{
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnvOrDefault("SMTP_USER", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", "hello@holger-hahn.net"),
		ToEmail:      getEnvOrDefault("TO_EMAIL", "hello@holger-hahn.net"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func sanitizeInput(input string) string {
	// Remove potential HTML/script content
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	return strings.TrimSpace(input)
}

func sendContactEmail(config *EmailConfig, req *ContactRequest) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.FromEmail)
	m.SetHeader("To", config.ToEmail)
	m.SetHeader("Subject", fmt.Sprintf("New Contact Form Submission from %s", req.Name))

	// Create HTML email body
	htmlBody := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<h2 style="color: #2563eb;">New Contact Form Submission</h2>
			<div style="background-color: #f8fafc; padding: 20px; border-radius: 8px; margin: 20px 0;">
				<p><strong>Name:</strong> %s</p>
				<p><strong>Company:</strong> %s</p>
				<p><strong>Email:</strong> <a href="mailto:%s">%s</a></p>
			</div>
			<div style="background-color: #ffffff; padding: 20px; border: 1px solid #e5e7eb; border-radius: 8px;">
				<h3 style="color: #374151; margin-top: 0;">Project Description:</h3>
				<p style="white-space: pre-wrap;">%s</p>
			</div>
			<hr style="margin: 30px 0; border: none; border-top: 1px solid #e5e7eb;">
			<p style="color: #6b7280; font-size: 14px;">
				This email was sent from the contact form on holger-hahn.net
			</p>
		</body>
		</html>
	`, 
		sanitizeInput(req.Name),
		sanitizeInput(req.Company),
		req.Email,
		req.Email,
		sanitizeInput(req.Project),
	)

	m.SetBody("text/html", htmlBody)

	// Create plain text version as fallback
	textBody := fmt.Sprintf(`
New Contact Form Submission

Name: %s
Company: %s
Email: %s

Project Description:
%s

---
This email was sent from the contact form on holger-hahn.net
`,
		req.Name,
		req.Company,
		req.Email,
		req.Project,
	)

	m.AddAlternative("text/plain", textBody)

	// Configure SMTP
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPassword)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func handleContactForm(config *EmailConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ContactRequest

		// Bind and validate form data
		if err := c.ShouldBind(&req); err != nil {
			log.Printf("Validation error: %v", err)
			c.JSON(http.StatusBadRequest, ContactResponse{
				Success: false,
				Message: "Please check all required fields and try again.",
			})
			return
		}

		// Additional email validation
		if !validateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, ContactResponse{
				Success: false,
				Message: "Please provide a valid email address.",
			})
			return
		}

		// Send email
		if err := sendContactEmail(config, &req); err != nil {
			log.Printf("Email sending error: %v", err)
			c.JSON(http.StatusInternalServerError, ContactResponse{
				Success: false,
				Message: "Sorry, there was an error sending your message. Please try again later or contact us directly.",
			})
			return
		}

		log.Printf("Contact form submitted successfully by %s (%s)", req.Name, req.Email)

		// Return success response
		c.JSON(http.StatusOK, ContactResponse{
			Success: true,
			Message: "Thank you for your message! I'll get back to you within 24 hours.",
		})
	}
}

func main() {
	// Initialize email configuration
	emailConfig := getEmailConfig()

	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Home page route
	r.GET("/", func(c *gin.Context) {
		component := templates.Index()
		c.Header("Content-Type", "text/html")
		component.Render(c.Request.Context(), c.Writer)
	})

	// Contact form submission route
	r.POST("/contact", handleContactForm(emailConfig))

	log.Println("Server starting on :8080")
	r.Run(":8080")
}