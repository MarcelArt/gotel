package web

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsLicenseFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"LICENSE", true},
		{"LICENCE", true},
		{"License", true},
		{"license.txt", true},
		{"license.md", true},
		{"notice", true},
		{"NOTICE.txt", true},
		{"copying", true},
		{"README.md", false},
		{"main.go", false},
		{"LICENSE.png", false},
	}

	for _, tt := range tests {
		result := isLicenseFile(tt.filename)
		if result != tt.expected {
			t.Errorf("isLicenseFile(%q) = %v; want %v", tt.filename, result, tt.expected)
		}
	}
}

func TestLoadLicenses(t *testing.T) {
	// Create a temporary mock license structure
	tmpDir, err := os.MkdirTemp("", "mock_licenses")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectories and mock license files
	pkg1 := filepath.Join(tmpDir, "github.com/test/pkg1")
	if err := os.MkdirAll(pkg1, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkg1, "LICENSE"), []byte("Mock MIT License for pkg1"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkg1, "NOTICE"), []byte("Mock Notice for pkg1"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	pkg2 := filepath.Join(tmpDir, "golang.org/x/pkg2")
	if err := os.MkdirAll(pkg2, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkg2, "License.txt"), []byte("Mock BSD License for pkg2"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	// Load the licenses
	licenses, err := LoadLicenses(tmpDir)
	if err != nil {
		t.Fatalf("LoadLicenses returned error: %v", err)
	}

	// Verify package count
	if len(licenses) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(licenses))
	}

	// Verify packages are sorted alphabetically
	if licenses[0].Package != "github.com/test/pkg1" {
		t.Errorf("expected first package to be github.com/test/pkg1, got %s", licenses[0].Package)
	}
	if licenses[1].Package != "golang.org/x/pkg2" {
		t.Errorf("expected second package to be golang.org/x/pkg2, got %s", licenses[1].Package)
	}

	// Verify pkg1 files
	if len(licenses[0].Files) != 2 {
		t.Errorf("expected 2 license files for pkg1, got %d", len(licenses[0].Files))
	}

	// Verify pkg2 files
	if len(licenses[1].Files) != 1 {
		t.Errorf("expected 1 license file for pkg2, got %d", len(licenses[1].Files))
	}
}
