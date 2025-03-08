package usecase

import (
	"context"
	"errors"
)

var (
	ErrPasswordRequired = errors.New("password is required for encrypted PDF")
)

type ConvertInput struct {
	FilePath string
	Password string
	Density  string
	Quality  int
}

type ConvertOutput struct {
	FileId        string
	EncodedImages []string
}

type ConvertUsecase struct {
	builder   FileBuilder
	converter ImageConvertService
	decrypter PdfDecryptService
}

func NewConvertUsecase(builder FileBuilder, converter ImageConvertService, decrypter PdfDecryptService) *ConvertUsecase {
	return &ConvertUsecase{
		builder:   builder,
		converter: converter,
		decrypter: decrypter,
	}
}

func (u *ConvertUsecase) Execute(ctx context.Context, input *ConvertInput) (*ConvertOutput, error) {
	file, err := u.builder.BuildFromPath(input.FilePath)
	if err != nil {
		return nil, err
	}

	isPasswordGiven := input.Password != ""
	if file.IsEncrypted() {
		if !isPasswordGiven {
			return nil, ErrPasswordRequired
		}

		if err = u.decrypter.Decrypt(ctx, file, input.Password); err != nil {
			return nil, err
		}
	}

	images, err := u.converter.Convert(ctx, file, ImageConvertOptions{
		Density: input.Density,
		Quality: input.Quality,
	})
	if err != nil {
		return nil, err
	}

	return &ConvertOutput{
		FileId:        file.Id(),
		EncodedImages: images,
	}, nil
}
