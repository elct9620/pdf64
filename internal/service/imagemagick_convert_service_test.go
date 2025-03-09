package service_test

import (
	"context"
	"encoding/base64"
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

	// Test cases for different conversion options
	testCases := []struct {
		name    string
		options usecase.ImageConvertOptions
	}{
		{
			name: "Standard conversion",
			options: usecase.ImageConvertOptions{
				Density: "150",
				Quality: 90,
				Merge:   false,
			},
		},
		{
			name: "Merged conversion",
			options: usecase.ImageConvertOptions{
				Density: "150",
				Quality: 90,
				Merge:   true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runConversionTest(t, file, service, tc.options)
		})
	}
}

func runConversionTest(t *testing.T, file *entity.File, service *service.ImageMagickConvertService, options usecase.ImageConvertOptions) {

	// Convert the PDF to images - this returns base64 encoded images
	base64Images, err := service.Convert(context.Background(), file, options)
	if err != nil {
		t.Fatalf("Failed to convert PDF to images: %v", err)
	}

	// Verify we got encoded images
	if len(base64Images) == 0 {
		t.Error("Expected at least one encoded image, got none")
	}

	// If merge is enabled, we should only get one image
	if options.Merge && len(base64Images) > 1 {
		t.Errorf("Expected only one merged image, got %d", len(base64Images))
	}

	// Check that the returned strings are valid base64 encoded data
	for i, encodedImage := range base64Images {
		if len(encodedImage) == 0 {
			t.Errorf("Encoded image %d is empty", i)
			continue
		}

		// Check if the string starts with the base64 image prefix
		if !strings.HasPrefix(encodedImage, "data:image/jpeg;base64,") {
			t.Errorf("Encoded image %d does not have valid image data prefix", i)
			continue
		}

		// Extract the base64 part
		base64Data := strings.TrimPrefix(encodedImage, "data:image/jpeg;base64,")

		// Try to decode it to verify it's valid base64
		_, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			t.Errorf("Encoded image %d contains invalid base64 data: %v", i, err)
			continue
		}

		t.Logf("Successfully verified base64 encoded image %d", i)
	}
}
