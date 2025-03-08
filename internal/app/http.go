package app

import (
	ctrlV1 "github.com/elct9620/pdf64/internal/controller/v1"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	chi.Router
}

func NewServer(
	ctrlV1 *ctrlV1.Service,
) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	v1.Register(r, ctrlV1)

	return &Server{
		Router: r,
	}
}
