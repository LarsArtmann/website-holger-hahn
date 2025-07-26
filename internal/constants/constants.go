// Package constants defines application-wide constants to avoid magic numbers
// and provide meaningful names for commonly used values.
package constants

import (
	"net/http"
	"os"
)

// File and directory permissions.
const (
	// DefaultDirectoryPerms is the default permission for created directories (rwxr-x---).
	DefaultDirectoryPerms os.FileMode = 0o750
)

// HTTP Status Codes.
const (
	// HTTPInternalServerError represents HTTP 500 Internal Server Error.
	HTTPInternalServerError = http.StatusInternalServerError
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
