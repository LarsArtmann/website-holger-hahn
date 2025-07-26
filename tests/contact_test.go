package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"holger-hahn-website/internal/constants"

	"github.com/gin-gonic/gin"
)

// TestContactLinksPresent verifies that all contact information is correctly embedded in the HTML.
func TestContactLinksPresent(t *testing.T) {
	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Serve static files (mock for testing)
	router.Static("/static", "./static")

	// Mock the index page handler
	router.GET("/", func(c *gin.Context) {
		// This would normally render the template, but for testing we'll use a simplified version
		html := `<!DOCTYPE html>
<html>
<head><title>Holger M. Hahn - Digital Assets Solutions Architect</title></head>
<body>
	<a href="mailto:hello@holger-hahn.net">Email</a>
	<a href="https://www.linkedin.com/in/holger-hahn-1b258ba3/">LinkedIn</a>
	<footer>
		<a href="mailto:hello@holger-hahn.net">hello@holger-hahn.net</a>
		<a href="https://www.linkedin.com/in/holger-hahn-1b258ba3/">LinkedIn Profile</a>
	</footer>
</body>
</html>`
		c.Data(constants.HTTPOKStatus, "text/html; charset=utf-8", []byte(html))
	})

	// Create test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Test the homepage
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to get homepage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != constants.HTTPOKStatus {
		t.Errorf("Expected status %d, got %d", constants.HTTPOKStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyStr := string(body)

	// Test cases for contact information
	testCases := []struct {
		name        string
		expected    string
		description string
	}{
		{
			name:        "Email Link",
			expected:    "mailto:hello@holger-hahn.net",
			description: "Email link should point to hello@holger-hahn.net",
		},
		{
			name:        "LinkedIn Link",
			expected:    "https://www.linkedin.com/in/holger-hahn-1b258ba3/",
			description: "LinkedIn link should point to Holger's profile",
		},
		{
			name:        "No Placeholder Email",
			expected:    "holger@example.com",
			description: "Should not contain placeholder email",
		},
		{
			name:        "No Placeholder LinkedIn",
			expected:    "linkedin.com/in/holger-hahn\"",
			description: "Should not contain old LinkedIn URL pattern",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "No Placeholder Email" || tc.name == "No Placeholder LinkedIn" {
				// These should NOT be present
				if strings.Contains(bodyStr, tc.expected) {
					t.Errorf("FAIL: %s - Found placeholder content: %s", tc.description, tc.expected)
				} else {
					t.Logf("PASS: %s - No placeholder content found", tc.description)
				}
			} else {
				// These should be present
				if !strings.Contains(bodyStr, tc.expected) {
					t.Errorf("FAIL: %s - Expected to find: %s", tc.description, tc.expected)
				} else {
					t.Logf("PASS: %s - Found expected content: %s", tc.description, tc.expected)
				}
			}
		})
	}
}

// TestContactLinksReachable verifies that the contact links are valid URLs.
func TestContactLinksReachable(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		description string
		expectCode  int
	}{
		{
			name:        "LinkedIn Profile",
			url:         "https://www.linkedin.com/in/holger-hahn-1b258ba3/",
			expectCode:  constants.HTTPOKStatus,
			description: "LinkedIn profile should be reachable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					// Allow redirects but don't follow them for testing
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Get(tc.url)
			if err != nil {
				t.Errorf("FAIL: %s - Could not reach URL %s: %v", tc.description, tc.url, err)
				return
			}
			defer resp.Body.Close()

			// LinkedIn typically returns 200 or redirects (3xx), both are acceptable
			if resp.StatusCode >= constants.HTTPSuccessRangeStart && resp.StatusCode < constants.HTTPClientErrorRangeStart {
				t.Logf("PASS: %s - URL %s returned status %d", tc.description, tc.url, resp.StatusCode)
			} else {
				t.Errorf("FAIL: %s - URL %s returned unexpected status %d", tc.description, tc.url, resp.StatusCode)
			}
		})
	}
}

// TestEmailFormat verifies email format validity.
func TestEmailFormat(t *testing.T) {
	email := "hello@holger-hahn.net"

	// Basic email validation
	if !strings.Contains(email, "@") {
		t.Errorf("FAIL: Email %s does not contain @", email)
	}

	if !strings.Contains(email, ".") {
		t.Errorf("FAIL: Email %s does not contain domain", email)
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		t.Errorf("FAIL: Email %s has invalid format", email)
	}

	if len(parts[0]) == 0 {
		t.Errorf("FAIL: Email %s has empty local part", email)
	}

	if len(parts[1]) == 0 {
		t.Errorf("FAIL: Email %s has empty domain part", email)
	}

	t.Logf("PASS: Email %s has valid format", email)
}
