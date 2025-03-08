package service_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/elct9620/pdf64/internal/entity"
	"github.com/elct9620/pdf64/internal/service"
	"github.com/elct9620/pdf64/internal/usecase"
)

func TestImageMagickConvertService_Convert(t *testing.T) {
	// Use the real PDF file from fixtures
	// Get the absolute path to the fixtures directory from project root
	// First find the project root directory
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	pdfPath := filepath.Join(projectRoot, "fixtures", "dummy.pdf")

	// Verify the test PDF file exists
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		t.Fatalf("Test PDF file not found at %s", pdfPath)
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

	// Convert the PDF to images - this returns base64 encoded images, not file paths
	encodedImages, err := service.Convert(context.Background(), file, options)
	if err != nil {
		t.Fatalf("Failed to convert PDF to images: %v", err)
	}

	// Verify we got encoded images
	if len(encodedImages) == 0 {
		t.Error("Expected at least one encoded image, got none")
	}

	// Check that the returned strings are valid base64 encoded data
	for i, encodedImage := range encodedImages {
		if len(encodedImage) == 0 {
			t.Errorf("Encoded image %d is empty", i)
			continue
		}
		
		// Check if the string starts with the base64 image prefix
		if !strings.HasPrefix(encodedImage, "data:image/") {
			t.Errorf("Encoded image %d does not have valid image data prefix", i)
		}
	}
}
