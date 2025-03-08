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
)

func TestApiV1Convert(t *testing.T) {
	// 使用表格驅動測試
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
				// UUID 格式驗證 (簡單檢查)
				if len(resp.Id) < 32 {
					t.Errorf("expected UUID format for Id, got: %s", resp.Id)
				}
				
				// 檢查 Data 陣列
				if len(resp.Data) != 0 {
					// 如果有數據，可以檢查每個項目是否為有效的 base64 編碼
					for i, item := range resp.Data {
						if item == "" {
							t.Errorf("Data[%d] is empty", i)
						}
					}
				}
			},
		},
		{
			name:           "無效品質參數測試",
			density:        "300",
			quality:        "invalid",
			fileContent:    "測試PDF內容",
			expectedStatus: http.StatusOK, // 目前實現會忽略無效品質，所以仍然返回 200
			validateResp: func(t *testing.T, resp *apiV1.ConvertResponse) {
				if resp.Id == "" {
					t.Errorf("expected non-empty ID, got empty string")
				}
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
			part, err := writer.CreateFormFile("data", "test.pdf")
			if err != nil {
				t.Fatal(err)
			}
			_, err = io.Copy(part, strings.NewReader(tt.fileContent))
			if err != nil {
				t.Fatal(err)
			}
			err = writer.Close()
			if err != nil {
				t.Fatal(err)
			}

			// 創建請求
			req := httptest.NewRequest("POST", "/v1/convert", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			
			// 執行請求
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			// 檢查狀態碼
			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, recorder.Code)
			}

			// 解析並驗證回應
			if tt.expectedStatus == http.StatusOK {
				var resp apiV1.ConvertResponse
				respBody := recorder.Body.Bytes()
				
				// 檢查回應是否為有效的 JSON
				if !json.Valid(respBody) {
					t.Fatalf("response is not valid JSON: %s", respBody)
				}
				
				err := json.Unmarshal(respBody, &resp)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				
				// 基本結構驗證
				if resp.Id == "" {
					t.Errorf("expected non-empty ID, got empty string")
				}
				
				if resp.Data == nil {
					t.Errorf("expected non-nil Data array, got nil")
				}
				
				// 檢查 Content-Type 標頭
				contentType := recorder.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected Content-Type to be application/json, got %s", contentType)
				}
				
				// 執行自定義驗證
				if tt.validateResp != nil {
					tt.validateResp(t, &resp)
				}
			}
		})
	}
}
