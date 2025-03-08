package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServiceImpl interface {
	Convert(ctx context.Context, req *ConvertRequest) (*ConvertResponse, error)
}

type Service struct {
	router chi.Router
}

func Register(r chi.Router, impl ServiceImpl) {
	r.Post("/v1/convert", PostConvert(impl))
}

func respondWithError(w http.ResponseWriter, err Error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
