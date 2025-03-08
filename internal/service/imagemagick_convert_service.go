package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/elct9620/pdf64/internal/entity"
	"github.com/elct9620/pdf64/internal/usecase"
)

// ImageMagickConvertService implements the ImageConvertService interface
// using ImageMagick's convert command
type ImageMagickConvertService struct{}

// NewImageMagickConvertService creates a new ImageMagickConvertService
func NewImageMagickConvertService() *ImageMagickConvertService {
	return &ImageMagickConvertService{}
}

// Convert converts a PDF file to base64 encoded images
func (s *ImageMagickConvertService) Convert(ctx context.Context, file *entity.File, options usecase.ImageConvertOptions) ([]string, error) {
	// Create temporary directory for output images
	tmpDir, err := os.MkdirTemp("", "pdf64-images-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Prepare output path pattern
	outputPattern := filepath.Join(tmpDir, "page-%d.jpg")

	// Prepare convert command arguments
	args := []string{
		"-density", options.Density,
		"-quality", fmt.Sprintf("%d", options.Quality),
		file.Path(),
		outputPattern,
	}

	// Execute convert command
	cmd := exec.CommandContext(ctx, "convert", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to convert PDF to images: %w", err)
	}

	// Read generated images and encode them to base64
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary directory: %w", err)
	}

	var encodedImages []string
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "page-") {
			imagePath := filepath.Join(tmpDir, f.Name())
			
			// Read image file
			imageData, err := os.ReadFile(imagePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read image file %s: %w", imagePath, err)
			}
			
			// Encode image to base64
			encoded := base64.StdEncoding.EncodeToString(imageData)
			encodedImages = append(encodedImages, encoded)
		}
	}

	return encodedImages, nil
}
