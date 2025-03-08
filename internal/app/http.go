package app

import (
	ctrlV1 "github.com/elct9620/pdf64/internal/controller/v1"
	v1 "github.com/elct9620/pdf64/pkg/apis/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

type Server struct {
	chi.Router
}

func NewServer(
	ctrlV1 *ctrlV1.Service,
) *Server {
	logger := httplog.NewLogger("pdf64", httplog.Options{
		JSON:    true,
		Concise: true,
	})

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	v1.Register(r, ctrlV1)

	return &Server{
		Router: r,
	}
}
