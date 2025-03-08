package v1

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/elct9620/pdf64/internal/usecase"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
)

func (s *Service) Convert(ctx context.Context, req *v1.ConvertRequest) (*v1.ConvertResponse, error) {
	// Create temporary file
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, "upload.pdf")
	
	// Create temporary file
	tmpFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()
	
	// Copy uploaded file content to temporary file
	_, err = io.Copy(tmpFile, req.File)
	if err != nil {
		return nil, err
	}
	
	// Ensure file is written and closed
	tmpFile.Close()
	
	// Delete temporary file when function exits
	defer os.Remove(filePath)
	
	// Parse quality parameter
	quality := 90 // Default quality
	if req.Quality > 0 {
		quality = req.Quality
	}

	// Use density from request or default
	density := "150" // Default density
	if req.Density != "" {
		density = req.Density
	}

	// Execute conversion use case
	out, err := s.convertUsecase.Execute(ctx, &usecase.ConvertInput{
		FilePath: filePath,
		Density:  density,
		Quality:  quality,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConvertResponse{
		Id:   out.FileId,
		Data: out.EncodedImages,
	}, nil
}
