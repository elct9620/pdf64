package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, err Error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

// Helper function to respond with JSON
func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

type ConvertRequest struct {
	Density string `json:"density"`
	Quality int    `json:"quality"`
	File    io.ReadCloser
}

type ConvertResponse struct {
	Id   string   `json:"id"`
	Data []string `json:"data"`
}

func PostConvert(impl ServiceImpl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form, 32 MB is the max memory
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			respondWithError(w, Error{
				Code:    ErrCodeBadRequest,
				Message: "Failed to parse form",
			}, http.StatusBadRequest)
			return
		}

		// Get optional form parameters
		density := r.FormValue("density")

		// Parse optional quality parameter
		quality := 0
		qualityStr := r.FormValue("quality")
		if qualityStr != "" {
			quality, err = strconv.Atoi(qualityStr)
			if err != nil {
				// Continue processing, let ServiceImpl handle validation
				quality = 0
			}
		}

		// Get uploaded file
		file, _, err := r.FormFile("data")
		if err != nil {
			respondWithError(w, Error{
				Code:    ErrCodeBadRequest,
				Message: "Failed to get uploaded file",
			}, http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 創建 ConvertRequest 物件
		req := ConvertRequest{
			Density: density,
			Quality: quality,
			File:    file,
		}

		// Call the implementation's Convert method
		resp, err := impl.Convert(r.Context(), &req)
		if err != nil {
			// Check if it's a custom error
			if apiErr, ok := err.(Error); ok {
				respondWithError(w, apiErr, http.StatusBadRequest)
			} else {
				// Unknown error
				respondWithError(w, Error{
					Code:    ErrCodeInternal,
					Message: "Conversion failed: " + err.Error(),
				}, http.StatusInternalServerError)
			}
			return
		}

		// Set response headers and return response
		respondWithJSON(w, resp, http.StatusOK)
	}
}
