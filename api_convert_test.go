package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elct9620/pdf64/internal/app"
	v1 "github.com/elct9620/pdf64/internal/controller/v1"
	apiV1 "github.com/elct9620/pdf64/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestApiV1Convert(t *testing.T) {
	tests := []struct {
		name           string
		density        string
		quality        string
		fileContent    string
		expectedStatus int
		validateResp   func(t *testing.T, resp *apiV1.ConvertResponse)
	}{
		{
			name:           "基本轉換測試",
			density:        "300",
			quality:        "90",
			fileContent:    "測試PDF內容",
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *apiV1.ConvertResponse) {
				assert.NotEmpty(t, resp.Id)
				assert.NotNil(t, resp.Data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 設置測試服務器
			apiV1Service := v1.NewService()
			server := app.NewServer(apiV1Service)

			// 創建一個測試請求
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			
			// 添加表單字段
			_ = writer.WriteField("density", tt.density)
			_ = writer.WriteField("quality", tt.quality)
			
			// 添加文件
			part, _ := writer.CreateFormFile("data", "test.pdf")
			io.Copy(part, strings.NewReader(tt.fileContent))
			writer.Close()

			// 創建請求
			req := httptest.NewRequest("POST", "/v1/convert", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			
			// 執行請求
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			// 檢查狀態碼
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// 解析並驗證回應
			if tt.expectedStatus == http.StatusOK {
				var resp apiV1.ConvertResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				
				if tt.validateResp != nil {
					tt.validateResp(t, &resp)
				}
			}
		})
	}
}
