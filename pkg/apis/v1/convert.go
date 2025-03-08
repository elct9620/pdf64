package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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
		// 解析多部分表單，32 MB 是最大記憶體
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "無法解析表單", http.StatusBadRequest)
			return
		}

		// 獲取可選的表單參數
		density := r.FormValue("density")
		
		// 解析可選的 quality 參數
		quality := 0
		qualityStr := r.FormValue("quality")
		if qualityStr != "" {
			quality, err = strconv.Atoi(qualityStr)
			if err != nil {
				http.Error(w, "quality 參數必須是整數", http.StatusBadRequest)
				return
			}
		}

		// 獲取上傳的檔案
		file, _, err := r.FormFile("data")
		if err != nil {
			http.Error(w, "無法獲取上傳的檔案", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 創建 ConvertRequest 物件
		req := ConvertRequest{
			Density: density,
			Quality: quality,
			File:    file,
		}

		// 呼叫實現的 Convert 方法
		resp, err := impl.Convert(r.Context(), req)
		if err != nil {
			http.Error(w, "轉換失敗: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 設置回應標頭
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// 編碼並回傳回應
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "無法編碼回應", http.StatusInternalServerError)
			return
		}
	}
}
