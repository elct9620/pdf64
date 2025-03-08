package usecase

import "github.com/elct9620/pdf64/internal/entity"

type FileBuilder interface {
	BuildFromPath(path string) (*entity.File, error)
}
