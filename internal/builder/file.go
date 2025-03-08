package builder

import (
	"os/exec"

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

	file := entity.NewFile(id.String(), path)

	// Check if the file is encrypted using qpdf
	if b.isEncrypted(path) {
		file.Encrypt()
	}

	return file, nil
}

// isEncrypted checks if a PDF file is encrypted using qpdf
func (b *FileBuilder) isEncrypted(path string) bool {
	cmd := exec.Command("qpdf", "--is-encrypted", path)
	err := cmd.Run()

	// Return code 0 means the file is encrypted
	return err == nil
}
