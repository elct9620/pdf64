package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/elct9620/pdf64/internal/entity"
)

// QpdfDecryptService implements the usecase.PdfDecryptService interface
type QpdfDecryptService struct{}

// NewQpdfDecryptService creates a new QpdfDecryptService instance
func NewQpdfDecryptService() *QpdfDecryptService {
	return &QpdfDecryptService{}
}

// Decrypt decrypts a PDF file using qpdf
func (s *QpdfDecryptService) Decrypt(ctx context.Context, file *entity.File, password string) error {
	// Create a temporary file for the decrypted output
	tmpDir := os.TempDir()
	decryptedPath := filepath.Join(tmpDir, fmt.Sprintf("decrypted-%s.pdf", file.Id()))

	// Run qpdf to decrypt the file
	cmd := exec.CommandContext(
		ctx,
		"qpdf",
		"--decrypt",
		"--password="+password,
		file.Path(),
		decryptedPath,
	)

	// Capture stderr for better error messages
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to decrypt PDF: %w, stderr: %s", err, stderr.String())
	}

	// Replace the original file with the decrypted one
	if err := os.Rename(decryptedPath, file.Path()); err != nil {
		return fmt.Errorf("failed to replace original file with decrypted file: %w", err)
	}

	// Mark the file as decrypted
	file.Decrypt()

	return nil
}
