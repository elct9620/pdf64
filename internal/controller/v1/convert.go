package v1

import (
	"context"

	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
)

func (s *Service) Convert(ctx context.Context, req *v1.ConvertRequest) (*v1.ConvertResponse, error) {
	return &v1.ConvertResponse{}, nil
}
