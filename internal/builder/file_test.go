package builder_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/elct9620/pdf64/internal/builder"
)

func TestFileBuilder_BuildFromPath(t *testing.T) {
	// Ensure qpdf is installed
	if _, err := exec.LookPath("qpdf"); err != nil {
		t.Fatalf("qpdf is required for testing: %v", err)
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "pdf64-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test PDF file
	unencryptedPath := filepath.Join(tmpDir, "unencrypted.pdf")
	if err := os.WriteFile(unencryptedPath, []byte("%PDF-1.5\n%%EOF\n"), 0644); err != nil {
		t.Fatalf("failed to create test PDF: %v", err)
	}

	// Create an encrypted PDF file using qpdf
	encryptedPath := filepath.Join(tmpDir, "encrypted.pdf")
	cmd := exec.Command("qpdf", "--encrypt", "password", "password", "40", "--", unencryptedPath, encryptedPath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create encrypted PDF: %v", err)
	}

	// Table-driven tests
	tests := []struct {
		name           string
		path           string
		expectedId     bool
		expectedPath   string
		expectedEncrypted bool
	}{
		{
			name:             "Basic File",
			path:             "/path/to/file.pdf",
			expectedId:       true,
			expectedPath:     "/path/to/file.pdf",
			expectedEncrypted: false,
		},
		{
			name:             "Unencrypted PDF",
			path:             unencryptedPath,
			expectedId:       true,
			expectedPath:     unencryptedPath,
			expectedEncrypted: false,
		},
		{
			name:             "Encrypted PDF",
			path:             encryptedPath,
			expectedId:       true,
			expectedPath:     encryptedPath,
			expectedEncrypted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileBuilder := builder.NewFileBuilder()
			file, err := fileBuilder.BuildFromPath(tt.path)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if file == nil {
				t.Fatal("expected file to not be nil")
			}

			if tt.expectedId && file.Id() == "" {
				t.Error("expected ID to not be empty")
			}

			if file.Path() != tt.expectedPath {
				t.Errorf("expected path to be %q, got %q", tt.expectedPath, file.Path())
			}

			if file.IsEncrypted() != tt.expectedEncrypted {
				t.Errorf("expected IsEncrypted() to be %v, got %v", tt.expectedEncrypted, file.IsEncrypted())
			}
		})
	}
}
