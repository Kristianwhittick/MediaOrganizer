package organizer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePaths(t *testing.T) {
	org := New(false)
	
	tests := []struct {
		name        string
		source      string
		output      string
		shouldError bool
	}{
		{"valid paths", "/tmp", "/tmp/output", false},
		{"path traversal in source", "/tmp/../etc", "/tmp/output", true},
		{"path traversal in output", "/tmp", "/tmp/../etc", true},
		{"non-existent source", "/nonexistent", "/tmp/output", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := org.validatePaths(tt.source, tt.output)
			if tt.shouldError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestIsMediaFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"photo.jpg", true},
		{"photo.JPEG", true},
		{"video.mp4", true},
		{"document.pdf", false},
		{"image.png", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := IsMediaFile(tt.filename)
			if result != tt.expected {
				t.Errorf("IsMediaFile(%s) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestGetUniqueFilename(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create a test file
	testFile := filepath.Join(tmpDir, "test.jpg")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}
	
	unique := GetUniqueFilename(tmpDir, "test.jpg")
	expected := filepath.Join(tmpDir, "test1.jpg")
	
	if unique != expected {
		t.Errorf("GetUniqueFilename() = %s, want %s", unique, expected)
	}
}