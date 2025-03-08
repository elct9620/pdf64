package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

// Convert converts a PDF file to images and returns their file paths
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

	// Collect paths of generated images
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary directory: %w", err)
	}

	var imagePaths []string
	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "page-") {
			imagePath := filepath.Join(tmpDir, f.Name())
			imagePaths = append(imagePaths, imagePath)
		}
	}

	// Sort the image paths to ensure correct page order
	// (This is important because ReadDir doesn't guarantee order)
	sort.Strings(imagePaths)

	return imagePaths, nil
}
