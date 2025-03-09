package v1

import (
	"io"
	"net/http"
	"strconv"
)

type ConvertRequest struct {
	Password string `json:"password"`
	Density  string `json:"density"`
	Quality  int    `json:"quality"`
	Merge    bool   `json:"merge"`
	File     io.ReadCloser
}

type ConvertResponse struct {
	Id   string   `json:"id"`
	Data []string `json:"data"`
}

func PostConvert(impl ServiceImpl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			respondWithError(w, r, Error{
				Code:    ErrCodeBadRequest,
				Message: "Failed to parse form",
			}, http.StatusBadRequest, err)
			return
		}

		density := r.FormValue("density")
		password := r.FormValue("password")

		quality := 0
		qualityStr := r.FormValue("quality")
		if qualityStr != "" {
			quality, err = strconv.Atoi(qualityStr)
			if err != nil {
				quality = 0
			}
		}

		merge := false
		mergeStr := r.FormValue("merge")
		if mergeStr == "true" || mergeStr == "1" {
			merge = true
		}

		file, _, err := r.FormFile("data")
		if err != nil {
			respondWithError(w, r, Error{
				Code:    ErrCodeBadRequest,
				Message: "Failed to get uploaded file",
			}, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		req := ConvertRequest{
			Password: password,
			Density:  density,
			Quality:  quality,
			Merge:    merge,
			File:     file,
		}

		resp, err := impl.Convert(r.Context(), &req)
		if err != nil {
			if apiErr, ok := err.(Error); ok {
				respondWithError(w, r, apiErr, http.StatusBadRequest, err)
			} else {
				respondWithError(w, r, Error{
					Code:    ErrCodeInternal,
					Message: "Conversion failed: " + err.Error(),
				}, http.StatusInternalServerError, err)
			}
			return
		}

		respondWithJSON(w, resp, http.StatusOK)
	}
}
