package builder

import (
	"github.com/elct9620/pdf64/internal/entity"
	"github.com/google/uuid"
)

// FileBuilder implements the usecase.FileBuilder interface
type FileBuilder struct{}

// NewFileBuilder creates a new FileBuilder instance
func NewFileBuilder() *FileBuilder {
	return &FileBuilder{}
}

// BuildFromPath creates a File entity from a file path
func (b *FileBuilder) BuildFromPath(path string) (*entity.File, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	
	return entity.NewFile(id.String(), path), nil
}
