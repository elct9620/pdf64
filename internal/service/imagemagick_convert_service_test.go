package service_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
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

	// Convert the PDF to images - this returns file paths to the generated images
	imagePaths, err := service.Convert(context.Background(), file, options)
	if err != nil {
		t.Fatalf("Failed to convert PDF to images: %v", err)
	}

	// Verify we got image paths
	if len(imagePaths) == 0 {
		t.Error("Expected at least one image path, got none")
	}

	// Create a temporary directory to copy the images to before they get deleted
	tmpDir, err := os.MkdirTemp("", "pdf64-test-images-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Check that the returned paths exist and are valid image files
	for i, path := range imagePaths {
		// Check if the file exists
		fileInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf("Image path %d does not exist: %s", i, path)
			continue
		}
		
		// Check if it's a file (not a directory)
		if fileInfo.IsDir() {
			t.Errorf("Image path %d is a directory, not a file: %s", i, path)
			continue
		}
		
		// Check if the file has content (size > 0)
		if fileInfo.Size() == 0 {
			t.Errorf("Image file %d is empty: %s", i, path)
			continue
		}
		
		// Copy the file to our temporary directory for further inspection if needed
		destPath := filepath.Join(tmpDir, filepath.Base(path))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read image file %d: %v", i, err)
			continue
		}
		
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			t.Errorf("Failed to copy image file %d: %v", i, err)
		}
		
		t.Logf("Successfully verified image %d: %s (copied to %s)", i, path, destPath)
	}
}
