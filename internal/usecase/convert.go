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
}

func NewConvertUsecase() *ConvertUsecase {
	return &ConvertUsecase{}
}

func (u *ConvertUsecase) Execute(ctx context.Context, input *ConvertInput) (*ConvertOutput, error) {
	return &ConvertOutput{
		FileId:        input.FileId,
		EncodedImages: []string{},
	}, nil
}
