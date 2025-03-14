package v1

import (
	"github.com/elct9620/pdf64/internal/usecase"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
)

var _ v1.ServiceImpl = &Service{}

type Service struct {
	convertUsecase *usecase.ConvertUsecase
}

func NewService(convertUsecase *usecase.ConvertUsecase) *Service {
	return &Service{
		convertUsecase: convertUsecase,
	}
}
