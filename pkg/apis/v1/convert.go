package v1

import (
	"io"
	"net/http"
)

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
	}
}
