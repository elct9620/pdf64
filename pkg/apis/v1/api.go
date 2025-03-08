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
	// 移除未使用的 router 字段
}

func Register(r chi.Router, impl ServiceImpl) {
	r.Post("/v1/convert", PostConvert(impl))
}

func respondWithError(w http.ResponseWriter, r *http.Request, err Error, statusCode int, originalErr error) {
	// 記錄錯誤到日誌
	if originalErr != nil {
		ctx := r.Context()
		httplog.LogEntrySetField(ctx, "error", slog.AnyValue(originalErr))
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		// 如果編碼失敗，記錄錯誤但不能再向客戶端發送響應
		ctx := r.Context()
		httplog.LogEntrySetField(ctx, "encode_error", slog.AnyValue(encodeErr))
	}
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// 如果編碼失敗，記錄錯誤但此時已無法向客戶端發送新的響應
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
