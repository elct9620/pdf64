package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/elct9620/pdf64/internal/entity"
	"github.com/elct9620/pdf64/internal/usecase"
	"github.com/go-chi/httplog/v2"
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
	}

	// If merge is enabled, we'll append all pages into a single image
	if options.Merge {
		// For merged output, we use a single file name
		outputPattern = filepath.Join(tmpDir, "merged.jpg")
		args = append(args, "-append", file.Path(), outputPattern)
	} else {
		args = append(args, file.Path(), outputPattern)
	}

	cliPath, err := exec.LookPath("magick")
	if err != nil {
		cliPath, err = exec.LookPath("convert")
		if err != nil {
			return nil, fmt.Errorf("failed to find ImageMagick command: %w", err)
		}
	}

	// Try to execute with magick command (ImageMagick 7)
	cmd := exec.CommandContext(ctx, cliPath, args...)

	// Capture stdout and stderr
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		httplog.LogEntrySetField(ctx, "stderr", slog.StringValue(stderr.String()))
		return nil, fmt.Errorf("failed to convert PDF to images: %w", err)
	}

	// Collect paths of generated images
	var imagePaths []string
	
	if options.Merge {
		// In merge mode, we only have one output file
		mergedPath := filepath.Join(tmpDir, "merged.jpg")
		if _, err := os.Stat(mergedPath); err == nil {
			imagePaths = append(imagePaths, mergedPath)
		} else {
			return nil, fmt.Errorf("failed to find merged output image: %w", err)
		}
	} else {
		// In normal mode, collect all page-* files
		files, err := os.ReadDir(tmpDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read temporary directory: %w", err)
		}

		for _, f := range files {
			if !f.IsDir() && strings.HasPrefix(f.Name(), "page-") {
				imagePath := filepath.Join(tmpDir, f.Name())
				imagePaths = append(imagePaths, imagePath)
			}
		}

		// Sort the image paths to ensure correct page order
		// (This is important because ReadDir doesn't guarantee order)
		sort.Strings(imagePaths)
	}

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
