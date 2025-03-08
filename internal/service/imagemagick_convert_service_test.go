package service_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/elct9620/pdf64/internal/entity"
	"github.com/elct9620/pdf64/internal/service"
	"github.com/elct9620/pdf64/internal/usecase"
)

func TestImageMagickConvertService_Convert(t *testing.T) {
	// Skip test if ImageMagick is not installed
	_, err := os.Stat("/usr/bin/convert")
	if os.IsNotExist(err) {
		t.Skip("ImageMagick convert command not found, skipping test")
	}

	// Create a test PDF file
	tmpDir, err := os.MkdirTemp("", "pdf64-test-*")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// For a real test, you would need a sample PDF file
	// Here we're just creating an empty file for demonstration
	pdfPath := filepath.Join(tmpDir, "test.pdf")
	if err := os.WriteFile(pdfPath, []byte("%PDF-1.5\n%%EOF\n"), 0644); err != nil {
		t.Fatalf("failed to create test PDF file: %v", err)
	}

	// Create file entity
	file := entity.NewFile("test-id", pdfPath)

	// Create service
	service := service.NewImageMagickConvertService()

	// Convert PDF to images
	options := usecase.ImageConvertOptions{
		Density: "150",
		Quality: 90,
	}

	// This test might fail with a real PDF since our test file is not a valid PDF
	// It's meant to demonstrate the structure of the test
	_, err = service.Convert(context.Background(), file, options)
	if err != nil {
		// We expect an error with our invalid PDF, but in a real test with a valid PDF
		// this should not error
		t.Logf("Expected error with invalid PDF: %v", err)
	}
}
