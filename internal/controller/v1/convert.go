package v1

import (
	"context"

	"github.com/elct9620/pdf64/internal/usecase"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
	"github.com/google/uuid"
)

func (s *Service) Convert(ctx context.Context, req *v1.ConvertRequest) (*v1.ConvertResponse, error) {
	convert := usecase.NewConvertUsecase()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	out, err := convert.Execute(ctx, &usecase.ConvertInput{
		FileId: id.String(),
	})

	return &v1.ConvertResponse{
		Id:   out.FileId,
		Data: out.EncodedImages,
	}, nil
}
