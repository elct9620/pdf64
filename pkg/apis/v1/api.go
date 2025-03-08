package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type ServiceImpl interface {
	Convert(ctx context.Context, req *ConvertRequest) (*ConvertResponse, error)
}

type Service struct {
}

func Register(r chi.Router, impl ServiceImpl) {
	r.Post("/v1/convert", PostConvert(impl))
}

func respondWithError(w http.ResponseWriter, r *http.Request, err Error, statusCode int, originalErr error) {
	if originalErr != nil {
		ctx := r.Context()
		httplog.LogEntrySetField(ctx, "error", slog.AnyValue(originalErr))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		ctx := r.Context()
		httplog.LogEntrySetField(ctx, "encode_error", slog.AnyValue(encodeErr))
	}
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
