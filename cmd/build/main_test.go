package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePath(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test_validate_path")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test file
	testFile := filepath.Join(subDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		path          string
		basePath      string
		isDestination bool
		expectError   bool
		expectedError error
	}{
		{
			name:        "valid file path within base",
			path:        testFile,
			basePath:    tmpDir,
			expectError: false,
		},
		{
			name:        "valid directory path within base",
			path:        subDir,
			basePath:    tmpDir,
			expectError: false,
		},
		{
			name:          "path traversal attempt source",
			path:          filepath.Join(tmpDir, "..", "outside"),
			basePath:      tmpDir,
			isDestination: false,
			expectError:   true,
			expectedError: ErrPathTraversal,
		},
		{
			name:          "path traversal attempt destination",
			path:          filepath.Join(tmpDir, "..", "outside"),
			basePath:      tmpDir,
			isDestination: true,
			expectError:   true,
			expectedError: ErrDestinationPathTraversal,
		},
		{
			name:        "base path equals path",
			path:        tmpDir,
			basePath:    tmpDir,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanPath, err := validatePath(tt.path, tt.basePath, tt.isDestination)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.expectedError != nil && err != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !filepath.IsAbs(cleanPath) {
					t.Errorf("expected absolute path, got %s", cleanPath)
				}
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	// Create temporary directories for testing
	srcDir, err := os.MkdirTemp("", "test_copy_src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "test_copy_dst")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	// Create source file
	srcFile := filepath.Join(srcDir, "test.txt")
	testContent := "test file content\nwith multiple lines"
	err = os.WriteFile(srcFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		srcPath     string
		dstPath     string
		expectError bool
	}{
		{
			name:        "copy valid file",
			srcPath:     srcFile,
			dstPath:     filepath.Join(dstDir, "copied.txt"),
			expectError: false,
		},
		{
			name:        "copy to nested directory",
			srcPath:     srcFile,
			dstPath:     filepath.Join(dstDir, "nested", "dir", "copied.txt"),
			expectError: false,
		},
		{
			name:        "copy nonexistent file",
			srcPath:     filepath.Join(srcDir, "nonexistent.txt"),
			dstPath:     filepath.Join(dstDir, "copied.txt"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyFile(tt.srcPath, tt.dstPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				// Verify file was copied correctly
				dstContent, err := os.ReadFile(tt.dstPath)
				if err != nil {
					t.Errorf("failed to read destination file: %v", err)
					return
				}

				if string(dstContent) != testContent {
					t.Errorf("file content mismatch. Expected %q, got %q", testContent, string(dstContent))
				}

				// Verify directory was created
				dstDir := filepath.Dir(tt.dstPath)
				if _, err := os.Stat(dstDir); os.IsNotExist(err) {
					t.Errorf("destination directory was not created: %s", dstDir)
				}
			}
		})
	}
}

func TestProcessFileOrDir(t *testing.T) {
	// Create temporary directories for testing
	srcDir, err := os.MkdirTemp("", "test_process_src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "test_process_dst")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	// Create test file
	testFile := filepath.Join(srcDir, "test.txt")
	testContent := "test file content"
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create test directory
	testSubDir := filepath.Join(srcDir, "subdir")
	err = os.MkdirAll(testSubDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		path        string
		isDir       bool
		expectError bool
	}{
		{
			name:        "process file",
			path:        testFile,
			isDir:       false,
			expectError: false,
		},
		{
			name:        "process directory",
			path:        testSubDir,
			isDir:       true,
			expectError: false,
		},
		{
			name:        "process source directory itself",
			path:        srcDir,
			isDir:       true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get file info
			info, err := os.Stat(tt.path)
			if err != nil {
				t.Fatal(err)
			}

			err = processFileOrDir(tt.path, srcDir, dstDir, info)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				// Calculate expected destination path
				relPath, err := filepath.Rel(srcDir, tt.path)
				if err != nil {
					t.Fatal(err)
				}
				expectedDstPath := filepath.Join(dstDir, relPath)

				// Verify destination exists and is correct type
				dstInfo, err := os.Stat(expectedDstPath)
				if err != nil {
					t.Errorf("destination was not created: %v", err)
					return
				}

				if dstInfo.IsDir() != tt.isDir {
					t.Errorf("destination type mismatch. Expected isDir=%v, got isDir=%v", tt.isDir, dstInfo.IsDir())
				}

				// If it's a file, verify content
				if !tt.isDir {
					dstContent, err := os.ReadFile(expectedDstPath)
					if err != nil {
						t.Errorf("failed to read destination file: %v", err)
						return
					}

					if string(dstContent) != testContent {
						t.Errorf("file content mismatch. Expected %q, got %q", testContent, string(dstContent))
					}
				}
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	// Create temporary directories for testing
	srcDir, err := os.MkdirTemp("", "test_copydir_src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "test_copydir_dst")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dstDir)

	// Create test structure
	testFiles := map[string]string{
		"file1.txt":           "content of file1",
		"file2.txt":           "content of file2",
		"subdir/file3.txt":    "content of file3",
		"subdir/file4.txt":    "content of file4",
		"subdir/sub2/file5.txt": "content of file5",
	}

	// Create test files and directories
	for filePath, content := range testFiles {
		fullPath := filepath.Join(srcDir, filePath)
		dir := filepath.Dir(fullPath)
		
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatal(err)
		}
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create empty directory
	emptyDir := filepath.Join(srcDir, "emptydir")
	err = os.MkdirAll(emptyDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		src         string
		dst         string
		expectError bool
	}{
		{
			name:        "copy entire directory",
			src:         srcDir,
			dst:         filepath.Join(dstDir, "copied"),
			expectError: false,
		},
		{
			name:        "copy nonexistent directory",
			src:         filepath.Join(srcDir, "nonexistent"),
			dst:         filepath.Join(dstDir, "copied"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyDir(tt.src, tt.dst)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				// Verify all files were copied correctly
				for filePath, expectedContent := range testFiles {
					srcPath := filepath.Join(tt.src, filePath)
					dstPath := filepath.Join(tt.dst, filePath)

					// Check source exists
					if _, err := os.Stat(srcPath); err != nil {
						continue // Skip if source doesn't exist for this test
					}

					// Check destination exists and has correct content
					dstContent, err := os.ReadFile(dstPath)
					if err != nil {
						t.Errorf("failed to read copied file %s: %v", dstPath, err)
						continue
					}

					if string(dstContent) != expectedContent {
						t.Errorf("file content mismatch for %s. Expected %q, got %q", 
							filePath, expectedContent, string(dstContent))
					}
				}

				// Verify empty directory was copied
				emptyDstDir := filepath.Join(tt.dst, "emptydir")
				if info, err := os.Stat(emptyDstDir); err != nil || !info.IsDir() {
					t.Errorf("empty directory was not copied correctly")
				}
			}
		})
	}
}

func TestStaticErrors(t *testing.T) {
	// Test that our static errors are properly defined
	if ErrPathTraversal == nil {
		t.Error("ErrPathTraversal should not be nil")
	}

	if ErrDestinationPathTraversal == nil {
		t.Error("ErrDestinationPathTraversal should not be nil")
	}

	// Test error messages
	if !strings.Contains(ErrPathTraversal.Error(), "path traversal") {
		t.Errorf("ErrPathTraversal should contain 'path traversal', got: %s", ErrPathTraversal.Error())
	}

	if !strings.Contains(ErrDestinationPathTraversal.Error(), "destination path traversal") {
		t.Errorf("ErrDestinationPathTraversal should contain 'destination path traversal', got: %s", ErrDestinationPathTraversal.Error())
	}
}