// Package security provides security utilities for safe file operations and input validation.
// This package contains reusable security functions to prevent common vulnerabilities
// like path traversal attacks and file inclusion vulnerabilities.
package security

import (
	"errors"
	"path/filepath"
	"strings"
)

// Static error variables to avoid dynamic error creation.
var (
	// ErrPathTraversal indicates a path traversal attempt was detected.
	ErrPathTraversal = errors.New("path traversal attempt detected")
	// ErrDestinationPathTraversal indicates a destination path traversal attempt was detected.
	ErrDestinationPathTraversal = errors.New("destination path traversal attempt detected")
)

// ValidatePath validates that a path is safe and within the expected base directory.
// It returns the cleaned absolute path if validation succeeds.
// This function prevents directory traversal attacks by ensuring paths stay within basePath.
//
// Parameters:
//   - path: The path to validate
//   - basePath: The base directory that path must be within
//   - isDestination: Whether this is a destination path (affects error message)
//
// Returns:
//   - string: The cleaned absolute path if validation succeeds
//   - error: Path traversal error if validation fails
func ValidatePath(path, basePath string, isDestination bool) (string, error) {
	// Clean and get absolute path to prevent path traversal.
	cleanPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	// Get clean base path for validation.
	cleanBasePath, err := filepath.Abs(filepath.Clean(basePath))
	if err != nil {
		return "", err
	}

	// Ensure the path is within the base directory (prevent directory traversal).
	if !strings.HasPrefix(cleanPath, cleanBasePath) {
		if isDestination {
			return "", ErrDestinationPathTraversal
		}

		return "", ErrPathTraversal
	}

	return cleanPath, nil
}

// ValidateFilePath validates a file path for safe file operations.
// This is a convenience wrapper around ValidatePath for source file paths.
func ValidateFilePath(path, basePath string) (string, error) {
	return ValidatePath(path, basePath, false)
}

// ValidateDestinationPath validates a destination path for safe file operations.
// This is a convenience wrapper around ValidatePath for destination file paths.
func ValidateDestinationPath(path, basePath string) (string, error) {
	return ValidatePath(path, basePath, true)
}
