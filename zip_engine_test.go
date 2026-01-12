package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCompressEngine(t *testing.T) {
	tests := []struct {
		name     string
		workers  int
		expected int
	}{
		{"default workers", 0, 4},
		{"negative workers", -1, 4},
		{"custom workers", 8, 8},
		{"single worker", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewCompressEngine(tt.workers)
			if engine.workers != tt.expected {
				t.Errorf("NewCompressEngine(%d) = %d workers, want %d", tt.workers, engine.workers, tt.expected)
			}
		})
	}
}

func TestCompressFileZip_SingleFile(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "zip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("Hello, World!")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create zip
	zipPath := filepath.Join(tempDir, "output.zip")
	engine := NewCompressEngine(2)
	if err := engine.CompressFileZip(tempDir, zipPath); err != nil {
		t.Fatalf("CompressFileZip failed: %v", err)
	}

	// Verify zip exists
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Fatal("zip file was not created")
	}

	// Verify zip contents
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer reader.Close()

	found := false
	for _, f := range reader.File {
		if f.Name == "test.txt" {
			found = true
			rc, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			buf := make([]byte, len(testContent))
			rc.Read(buf)
			rc.Close()
			if string(buf) != string(testContent) {
				t.Errorf("content mismatch: got %s, want %s", buf, testContent)
			}
		}
	}
	if !found {
		t.Error("test.txt not found in zip")
	}
}

func TestCompressFileZip_Directory(t *testing.T) {
	// Create temp directory structure
	tempDir, err := os.MkdirTemp("", "zip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create files
	files := map[string]string{
		"root.txt":        "root content",
		"subdir/file.txt": "subdir content",
	}
	for name, content := range files {
		path := filepath.Join(tempDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create zip
	zipPath := filepath.Join(tempDir, "output.zip")
	engine := NewCompressEngine(4)
	if err := engine.CompressFileZip(tempDir, zipPath); err != nil {
		t.Fatalf("CompressFileZip failed: %v", err)
	}

	// Verify zip contents
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer reader.Close()

	expectedFiles := map[string]bool{
		"root.txt":        false,
		"subdir/":         false,
		"subdir/file.txt": false,
	}

	for _, f := range reader.File {
		// Normalize path separators for Windows
		name := filepath.ToSlash(f.Name)
		if _, ok := expectedFiles[name]; ok {
			expectedFiles[name] = true
		}
	}

	for name, found := range expectedFiles {
		if !found {
			t.Errorf("expected file %s not found in zip", name)
		}
	}
}

func TestCompressFileZip_EmptyDirectory(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "zip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create source directory (separate from where zip will be created)
	sourceDir := filepath.Join(tempDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create zip of empty directory
	zipPath := filepath.Join(tempDir, "output.zip")
	engine := NewCompressEngine(2)
	if err := engine.CompressFileZip(sourceDir, zipPath); err != nil {
		t.Fatalf("CompressFileZip failed: %v", err)
	}

	// Verify zip exists and is valid
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	defer reader.Close()

	if len(reader.File) != 0 {
		t.Errorf("expected empty zip, got %d files", len(reader.File))
	}
}

func TestCompressFileZip_InvalidSource(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "output.zip")
	engine := NewCompressEngine(2)

	err = engine.CompressFileZip("/nonexistent/path", zipPath)
	if err == nil {
		t.Error("expected error for nonexistent source path")
	}
}

func TestCompressFileZip_InvalidDest(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	engine := NewCompressEngine(2)
	err = engine.CompressFileZip(tempDir, "/nonexistent/path/output.zip")
	if err == nil {
		t.Error("expected error for invalid destination path")
	}
}

func TestCompressFileZip_LargeFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zip_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a larger file (1MB)
	largeFile := filepath.Join(tempDir, "large.bin")
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	if err := os.WriteFile(largeFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(tempDir, "output.zip")
	engine := NewCompressEngine(4)
	if err := engine.CompressFileZip(tempDir, zipPath); err != nil {
		t.Fatalf("CompressFileZip failed: %v", err)
	}

	// Verify zip was created
	info, err := os.Stat(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("zip file is empty")
	}
}
