// Package main provides a build utility for generating static website files.
// It compiles templates, copies static assets, and creates the public directory
// structure for deployment.
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"holger-hahn-website/internal/constants"
	"holger-hahn-website/internal/security"
	"holger-hahn-website/templates"
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

// copyFile copies a single file from cleanSrcPath to cleanDstPath.
// Paths must already be validated for security.
func copyFile(cleanSrcPath, cleanDstPath string) error {
	// Open source file.
	srcFile, err := os.Open(cleanSrcPath) // #nosec G304 - Path validated by validatePath
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Ensure destination directory exists.
	err = os.MkdirAll(filepath.Dir(cleanDstPath), constants.DefaultDirectoryPerms)
	if err != nil {
		return err
	}

	// Create destination file.
	dstFile, err := os.Create(cleanDstPath) // #nosec G304 - Path validated by validatePath
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy file contents.
	_, err = dstFile.ReadFrom(srcFile)
	return err
}

// processFileOrDir processes a single file or directory entry during directory copying.
// It validates paths and either creates a directory or copies a file.
func processFileOrDir(path, src, dst string, info os.FileInfo) error {
	// Validate source path.
	cleanPath, err := validatePath(path, src, false)
	if err != nil {
		return err
	}

	// Calculate the destination path.
	relPath, err := filepath.Rel(src, path)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(dst, relPath)

	// Validate destination path.
	cleanDstPath, err := validatePath(dstPath, dst, true)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return os.MkdirAll(cleanDstPath, info.Mode())
	}

	// Copy file using the extracted copyFile function.
	return copyFile(cleanPath, cleanDstPath)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return processFileOrDir(path, src, dst, info)
	})
}
