package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"holger-hahn-website/templates"
)

func main() {
	// Create public directory for Firebase hosting
	publicDir := "public"
	err := os.MkdirAll(publicDir, 0755)
	if err != nil {
		panic(err)
	}

	// Generate index.html
	fmt.Println("Generating index.html...")
	indexFile, err := os.Create(filepath.Join(publicDir, "index.html"))
	if err != nil {
		panic(err)
	}
	defer indexFile.Close()

	// Render the index template
	component := templates.Index()
	err = component.Render(context.Background(), indexFile)
	if err != nil {
		panic(err)
	}

	// Copy static assets
	fmt.Println("Copying static assets...")
	err = copyDir("static", filepath.Join(publicDir, "static"))
	if err != nil {
		panic(err)
	}

	fmt.Println("Build completed! Files generated in public/ directory")
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate the destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = dstFile.ReadFrom(srcFile)
		return err
	})
}