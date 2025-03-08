package v1

import v1 "github.com/elct9620/pdf64/pkg/apis/v1"

var _ v1.ServiceImpl = &Service{}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
