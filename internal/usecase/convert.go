package usecase

import "context"

type ConvertInput struct {
	FileId   string
	FilePath string
}

type ConvertOutput struct {
	FileId        string
	EncodedImages []string
}

type ConvertUsecase struct {
	builder FileBuilder
}

func NewConvertUsecase(builder FileBuilder) *ConvertUsecase {
	return &ConvertUsecase{
		builder: builder,
	}
}

func (u *ConvertUsecase) Execute(ctx context.Context, input *ConvertInput) (*ConvertOutput, error) {
	file, err := u.builder.BuildFromPath(input.FilePath)
	if err != nil {
		return nil, err
	}

	return &ConvertOutput{
		FileId:        file.Id(),
		EncodedImages: []string{},
	}, nil
}
