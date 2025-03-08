package v1

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/elct9620/pdf64/internal/usecase"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
	"github.com/google/uuid"
)

func (s *Service) Convert(ctx context.Context, req *v1.ConvertRequest) (*v1.ConvertResponse, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	// Create temporary file
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, id.String()+".pdf")
	
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
	
	// Execute conversion use case
	out, err := s.convertUsecase.Execute(ctx, &usecase.ConvertInput{
		FileId:   id.String(),
		FilePath: filePath,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ConvertResponse{
		Id:   out.FileId,
		Data: out.EncodedImages,
	}, nil
}
