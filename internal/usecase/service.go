package usecase

import (
	"context"

	"github.com/elct9620/pdf64/internal/entity"
)

type ImageConvertOptions struct {
	Density string
	Quality int
	Merge   bool
}

type ImageConvertService interface {
	Convert(ctx context.Context, file *entity.File, options ImageConvertOptions) ([]string, error)
}

type PdfDecryptService interface {
	Decrypt(ctx context.Context, file *entity.File, password string) error
}
