package service

import (
	"context"
	"encoding/base64"
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

// Convert converts a PDF file to images and returns them as base64 encoded strings
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

	// Execute with magick command (ImageMagick 7)
	cmd := exec.CommandContext(ctx, "magick", args...)
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

	// Convert images to base64
	var base64Images []string
	for _, imagePath := range imagePaths {
		// Read image file
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file: %w", err)
		}

		// Encode to base64
		base64Data := base64.StdEncoding.EncodeToString(imageData)
		
		// Add data URI prefix
		base64Image := fmt.Sprintf("data:image/jpeg;base64,%s", base64Data)
		
		base64Images = append(base64Images, base64Image)
	}

	return base64Images, nil
}
