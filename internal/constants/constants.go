// Package constants defines application-wide constants to avoid magic numbers
// and provide meaningful names for commonly used values.
package constants

import (
	"net/http"
	"os"
	"time"
)

// File and directory permissions.
const (
	// DefaultDirectoryPerms is the default permission for created directories (rwxr-x---).
	DefaultDirectoryPerms os.FileMode = 0o750
)

// HTTP Status Codes.
const (
	// HTTPOKStatus represents HTTP 200 OK.
	HTTPOKStatus = http.StatusOK

	// HTTPBadRequestStatus represents HTTP 400 Bad Request.
	HTTPBadRequestStatus = http.StatusBadRequest

	// HTTPInternalServerError represents HTTP 500 Internal Server Error.
	HTTPInternalServerError = http.StatusInternalServerError
)

// Network Port Limits.
const (
	// MinValidPort is the minimum valid network port number.
	MinValidPort = 1

	// MaxValidPort is the maximum valid network port number.
	MaxValidPort = 65535
)

// Server Configuration Defaults.
const (
	// DefaultServerPort is the default port the server listens on.
	DefaultServerPort = 8080

	// DefaultReadTimeoutSeconds is the default HTTP read timeout in seconds.
	DefaultReadTimeoutSeconds = 30

	// DefaultWriteTimeoutSeconds is the default HTTP write timeout in seconds.
	DefaultWriteTimeoutSeconds = 30
)

// Database Configuration Defaults.
const (
	// DefaultMaxOpenConnections is the default maximum number of open database connections.
	DefaultMaxOpenConnections = 10

	// DefaultMaxIdleConnections is the default maximum number of idle database connections.
	DefaultMaxIdleConnections = 5
)

// Exit Status Codes.
const (
	// ExitFailure represents a non-zero exit status code for failures.
	ExitFailure = 1
)

// Performance Thresholds.
const (
	// FastResponseThreshold is the threshold under which responses are considered fast.
	FastResponseThreshold = 100 * time.Millisecond

	// AcceptableResponseThreshold is the threshold under which responses are acceptable.
	AcceptableResponseThreshold = 200 * time.Millisecond

	// SlowResponseThreshold is the threshold under which responses are slow but tolerable.
	SlowResponseThreshold = 500 * time.Millisecond

	// UnacceptableResponseThreshold is the threshold above which responses are unacceptable.
	UnacceptableResponseThreshold = 1 * time.Second
)

// HTTP Status Code Ranges.
const (
	// HTTPSuccessRangeStart is the start of 2xx HTTP status codes.
	HTTPSuccessRangeStart = 200

	// HTTPClientErrorRangeStart is the start of 4xx HTTP status codes.
	HTTPClientErrorRangeStart = 400
)

// Application Configuration.
const (
	// AppVersion is the current application version.
	AppVersion = "1.0.0"
)
