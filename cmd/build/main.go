// Package main provides a build utility for generating static website files.
// It compiles templates, copies static assets, and creates the public directory
// structure for deployment.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"holger-hahn-website/internal/constants"
	"holger-hahn-website/templates"
)

// Static error variables to avoid dynamic error creation.
var (
	// ErrPathTraversal indicates a path traversal attempt was detected.
	ErrPathTraversal = errors.New("path traversal attempt detected")
	// ErrDestinationPathTraversal indicates a destination path traversal attempt was detected.
	ErrDestinationPathTraversal = errors.New("destination path traversal attempt detected")
)

func main() {
	// Set up structured logging with build prefix.
	log.SetPrefix("[BUILD] ")
	// Create public directory for Firebase hosting.
	publicDir := "public"

	err := os.MkdirAll(publicDir, constants.DefaultDirectoryPerms)
	if err != nil {
		panic(err)
	}

	// Generate TailwindCSS first.
	log.Printf("Generating TailwindCSS...")

	err = os.Chdir(".")
	if err != nil {
		panic(err)
	}

	// Generate index.html.
	log.Printf("Generating index.html...")

	indexFile, err := os.Create(filepath.Join(publicDir, "index.html"))
	if err != nil {
		panic(err)
	}

	defer indexFile.Close()

	// Render the index template.
	component := templates.Index()

	err = component.Render(context.Background(), indexFile)
	if err != nil {
		panic(err)
	}

	// Copy static assets.
	log.Printf("Copying static assets...")

	err = copyDir("static", filepath.Join(publicDir, "static"))
	if err != nil {
		panic(err)
	}

	log.Printf("Build completed! Files generated in public/ directory")
}

func copyDir(src, dst string) error {
	// Clean and get absolute path for source directory to prevent path traversal.
	cleanSrc, err := filepath.Abs(filepath.Clean(src))
	if err != nil {
		return err
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Clean the path and validate it's within the source directory.
		cleanPath, err := filepath.Abs(filepath.Clean(path))
		if err != nil {
			return err
		}

		// Ensure the path is within the source directory (prevent directory traversal).
		if !strings.HasPrefix(cleanPath, cleanSrc) {
			return fmt.Errorf("path traversal attempt detected: %s", path)
		}

		// Calculate the destination path.
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		// Clean and validate destination path to prevent directory traversal.
		cleanDstPath, err := filepath.Abs(filepath.Clean(dstPath))
		if err != nil {
			return err
		}

		// Get clean destination directory for validation.
		cleanDst, err := filepath.Abs(filepath.Clean(dst))
		if err != nil {
			return err
		}

		// Ensure destination path is within the target directory.
		if !strings.HasPrefix(cleanDstPath, cleanDst) {
			return fmt.Errorf("destination path traversal attempt detected: %s", dstPath)
		}

		if info.IsDir() {
			return os.MkdirAll(cleanDstPath, info.Mode())
		}

		// Copy file - now using validated cleanPath.
		srcFile, err := os.Open(cleanPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		err = os.MkdirAll(filepath.Dir(cleanDstPath), constants.DefaultDirectoryPerms)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(cleanDstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = dstFile.ReadFrom(srcFile)

		return err
	})
}
