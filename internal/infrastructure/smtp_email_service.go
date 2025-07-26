package infrastructure

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"holger-hahn-website/internal/domain"

	"gopkg.in/gomail.v2"
)

// SMTPEmailService implements EmailService using SMTP.
type SMTPEmailService struct {
	dialer   *gomail.Dialer
	fromAddr string
	toAddr   string
}

// NewSMTPEmailService creates a new SMTP email service.
func NewSMTPEmailService() *SMTPEmailService {
	host := getEnvOrDefault("SMTP_HOST", "localhost")
	portStr := getEnvOrDefault("SMTP_PORT", "587")
	username := getEnvOrDefault("SMTP_USERNAME", "")
	password := getEnvOrDefault("SMTP_PASSWORD", "")
	fromAddr := getEnvOrDefault("FROM_EMAIL", "hello@holger-hahn.net")
	toAddr := getEnvOrDefault("TO_EMAIL", "hello@holger-hahn.net")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587 // Default SMTP port
	}

	dialer := gomail.NewDialer(host, port, username, password)

	// For development, you can disable TLS
	if getEnvOrDefault("SMTP_TLS", "true") == "false" {
		dialer.TLSConfig = nil
	}

	return &SMTPEmailService{
		dialer:   dialer,
		fromAddr: fromAddr,
		toAddr:   toAddr,
	}
}

// SendContactNotification sends a notification email about a new contact.
func (s *SMTPEmailService) SendContactNotification(ctx context.Context, contact *domain.Contact) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromAddr)
	m.SetHeader("To", s.toAddr)
	m.SetHeader("Subject", "New Contact Form Submission - "+contact.Name)

	body := fmt.Sprintf(`
New contact form submission received:

Name: %s
Company: %s
Email: %s
Project Description:
%s

Submitted at: %s
Contact ID: %s

Please respond within 24 hours.
`,
		contact.Name,
		contact.Company,
		contact.Email,
		contact.Project,
		contact.SubmittedAt.Format("2006-01-02 15:04:05 UTC"),
		contact.ID,
	)

	m.SetBody("text/plain", body)

	// In development mode, just log instead of sending
	if getEnvOrDefault("EMAIL_MODE", "development") == "development" {
		fmt.Printf("\n=== EMAIL NOTIFICATION ===\n")
		fmt.Printf("To: %s\n", s.toAddr)
		fmt.Printf("Subject: %s\n", "New Contact Form Submission - "+contact.Name)
		fmt.Printf("Body:\n%s\n", body)
		fmt.Printf("=========================\n\n")

		return nil
	}

	return s.dialer.DialAndSend(m)
}

// SendConfirmationEmail sends a confirmation email to the contact.
func (s *SMTPEmailService) SendConfirmationEmail(ctx context.Context, contact *domain.Contact) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromAddr)
	m.SetHeader("To", contact.Email)
	m.SetHeader("Subject", "Thank you for contacting Holger M. Hahn")

	body := fmt.Sprintf(`
Dear %s,

Thank you for reaching out! I've received your message about your digital asset project.

Here's a summary of your submission:
- Project: %s
- Submitted: %s

I'll review your requirements and get back to you within 24 hours with next steps.

Best regards,
Holger M. Hahn
Digital Assets Solutions Architect

---
This is an automated confirmation. Please don't reply to this email.
For urgent matters, contact me directly at hello@holger-hahn.net
`,
		contact.Name,
		contact.Project,
		contact.SubmittedAt.Format("January 2, 2006"),
	)

	m.SetBody("text/plain", body)

	// In development mode, just log instead of sending
	if getEnvOrDefault("EMAIL_MODE", "development") == "development" {
		fmt.Printf("\n=== CONFIRMATION EMAIL ===\n")
		fmt.Printf("To: %s\n", contact.Email)
		fmt.Printf("Subject: %s\n", "Thank you for contacting Holger M. Hahn")
		fmt.Printf("Body:\n%s\n", body)
		fmt.Printf("===========================\n\n")

		return nil
	}

	return s.dialer.DialAndSend(m)
}

// getEnvOrDefault gets an environment variable or returns a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
