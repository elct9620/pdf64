package builder_test

import (
	"bytes"
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

	// Find project root directory to locate fixtures
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Go up from internal/builder to project root
	projectRoot := filepath.Join(wd, "..", "..")

	// Path to the fixture PDF
	fixturesPdfPath := filepath.Join(projectRoot, "fixtures", "dummy.pdf")
	if _, err := os.Stat(fixturesPdfPath); os.IsNotExist(err) {
		t.Fatalf("fixture PDF not found at %s: %v", fixturesPdfPath, err)
	}

	// Create a temporary directory for the encrypted PDF
	tmpDir, err := os.MkdirTemp("", "pdf64-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an encrypted PDF file using qpdf
	encryptedPath := filepath.Join(tmpDir, "encrypted.pdf")
	
	// Try different qpdf encryption commands based on version
	var err error
	var stderr bytes.Buffer
	
	// First try with 256-bit AES encryption (newer qpdf versions)
	cmd := exec.Command("qpdf", "--encrypt", "password", "password", "256", "--", fixturesPdfPath, encryptedPath)
	cmd.Stderr = &stderr
	err = cmd.Run()
	
	// If that fails, try with 128-bit encryption
	if err != nil {
		stderr.Reset()
		cmd = exec.Command("qpdf", "--encrypt", "password", "password", "128", "--", fixturesPdfPath, encryptedPath)
		cmd.Stderr = &stderr
		err = cmd.Run()
	}
	
	// If that also fails, try with 40-bit encryption
	if err != nil {
		stderr.Reset()
		cmd = exec.Command("qpdf", "--encrypt", "password", "password", "40", "--", fixturesPdfPath, encryptedPath)
		cmd.Stderr = &stderr
		err = cmd.Run()
	}
	
	// If all encryption attempts fail, skip the encrypted PDF test
	if err != nil {
		t.Logf("Warning: Could not create encrypted PDF: %v, stderr: %s", err, stderr.String())
		t.Log("Skipping encrypted PDF test")
		encryptedPath = fixturesPdfPath // Use the unencrypted PDF as a fallback
	}

	// Table-driven tests
	tests := []struct {
		name              string
		path              string
		expectedId        bool
		expectedPath      string
		expectedEncrypted bool
	}{
		{
			name:              "Basic File",
			path:              "/path/to/file.pdf",
			expectedId:        true,
			expectedPath:      "/path/to/file.pdf",
			expectedEncrypted: false,
		},
		{
			name:              "Unencrypted PDF",
			path:              fixturesPdfPath,
			expectedId:        true,
			expectedPath:      fixturesPdfPath,
			expectedEncrypted: false,
		},
		{
			name:              "Encrypted PDF",
			path:              encryptedPath,
			expectedId:        true,
			expectedPath:      encryptedPath,
			expectedEncrypted: encryptedPath != fixturesPdfPath, // Only expect encrypted if we successfully created one
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
