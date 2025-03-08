package usecase

import "context"

type ConvertInput struct {
	FilePath string
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
}

func NewConvertUsecase(builder FileBuilder, converter ImageConvertService) *ConvertUsecase {
	return &ConvertUsecase{
		builder:   builder,
		converter: converter,
	}
}

func (u *ConvertUsecase) Execute(ctx context.Context, input *ConvertInput) (*ConvertOutput, error) {
	file, err := u.builder.BuildFromPath(input.FilePath)
	if err != nil {
		return nil, err
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
