package service_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/elct9620/pdf64/internal/entity"
	"github.com/elct9620/pdf64/internal/service"
)

func TestQpdfDecryptService_Decrypt(t *testing.T) {
	// Ensure qpdf is installed
	if _, err := exec.LookPath("qpdf"); err != nil {
		t.Fatalf("qpdf is required for testing: %v", err)
	}

	// Find project root directory to locate fixtures
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	
	// Go up from internal/service to project root
	projectRoot := filepath.Join(wd, "..", "..")
	
	// Path to the fixture PDF
	fixturesPdfPath := filepath.Join(projectRoot, "fixtures", "dummy.pdf")
	if _, err := os.Stat(fixturesPdfPath); os.IsNotExist(err) {
		t.Fatalf("fixture PDF not found at %s: %v", fixturesPdfPath, err)
	}
	
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "pdf64-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Copy the fixture PDF to the temp directory
	encryptedPath := filepath.Join(tmpDir, "encrypted.pdf")
	encryptedContent, err := os.ReadFile(fixturesPdfPath)
	if err != nil {
		t.Fatalf("failed to read fixture PDF: %v", err)
	}
	if err := os.WriteFile(encryptedPath, encryptedContent, 0644); err != nil {
		t.Fatalf("failed to write encrypted PDF: %v", err)
	}
	
	// Create an encrypted PDF file using qpdf
	password := "testpassword"
	var stderr bytes.Buffer
	cmd := exec.Command("qpdf", "--encrypt", password, password, "256", "--", encryptedPath, encryptedPath+".tmp")
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create encrypted PDF: %v, stderr: %s", err, stderr.String())
	}
	
	// Replace the original file with the encrypted one
	if err := os.Rename(encryptedPath+".tmp", encryptedPath); err != nil {
		t.Fatalf("failed to replace original file with encrypted file: %v", err)
	}
	
	// Create a file entity
	file := entity.NewFile("test-id", encryptedPath)
	file.Encrypt() // Mark as encrypted
	
	// Create the decrypt service
	decryptService := service.NewQpdfDecryptService()
	
	// Test with correct password
	err = decryptService.Decrypt(context.Background(), file, password)
	if err != nil {
		t.Errorf("failed to decrypt PDF with correct password: %v", err)
	}
	
	if file.IsEncrypted() {
		t.Error("file should be marked as decrypted after successful decryption")
	}
	
	// Verify the file is actually decrypted by checking if qpdf reports it as encrypted
	cmd = exec.Command("qpdf", "--is-encrypted", encryptedPath)
	if err := cmd.Run(); err == nil {
		t.Error("file should not be encrypted after decryption")
	}
	
	// Test with incorrect password (re-encrypt the file first)
	stderr.Reset()
	cmd = exec.Command("qpdf", "--encrypt", password, password, "256", "--", encryptedPath, encryptedPath+".tmp")
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to re-encrypt PDF: %v, stderr: %s", err, stderr.String())
	}
	
	if err := os.Rename(encryptedPath+".tmp", encryptedPath); err != nil {
		t.Fatalf("failed to replace original file with re-encrypted file: %v", err)
	}
	
	file.Encrypt() // Mark as encrypted again
	
	// Test with incorrect password
	err = decryptService.Decrypt(context.Background(), file, "wrongpassword")
	if err == nil {
		t.Error("decryption should fail with incorrect password")
	}
}
