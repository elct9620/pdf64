package builder_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/elct9620/pdf64/internal/builder"
)

func TestFileBuilder_BuildFromPath(t *testing.T) {
	// Arrange
	fileBuilder := builder.NewFileBuilder()
	path := "/path/to/file.pdf"

	// Act
	file, err := fileBuilder.BuildFromPath(path)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if file == nil {
		t.Fatal("expected file to not be nil")
	}

	if file.Path() != path {
		t.Errorf("expected path to be %q, got %q", path, file.Path())
	}

	if file.Id() == "" {
		t.Error("expected ID to not be empty")
	}
}

func TestFileBuilder_BuildFromPath_WithEncryptedFile(t *testing.T) {
	// Skip if qpdf is not installed
	if _, err := exec.LookPath("qpdf"); err != nil {
		t.Skip("qpdf not installed, skipping test")
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
	cmd := exec.Command("qpdf", "--encrypt", "user", "owner", "128", "--", unencryptedPath, encryptedPath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create encrypted PDF: %v", err)
	}

	// Test with unencrypted file
	fileBuilder := builder.NewFileBuilder()
	
	unencryptedFile, err := fileBuilder.BuildFromPath(unencryptedPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if unencryptedFile.IsEncrypted() {
		t.Error("unencrypted file should not be marked as encrypted")
	}

	// Test with encrypted file
	encryptedFile, err := fileBuilder.BuildFromPath(encryptedPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !encryptedFile.IsEncrypted() {
		t.Error("encrypted file should be marked as encrypted")
	}
}
