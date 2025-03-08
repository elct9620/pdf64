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
	"github.com/elct9620/pdf64/internal/builder"
	v1 "github.com/elct9620/pdf64/internal/controller/v1"
	"github.com/elct9620/pdf64/internal/service"
	"github.com/elct9620/pdf64/internal/usecase"
	apiV1 "github.com/elct9620/pdf64/pkg/apis/v1"
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
			name:           "Basic Conversion Test",
			density:        "300",
			quality:        "90",
			fileContent:    "Test PDF Content",
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *apiV1.ConvertResponse) {
				if len(resp.Id) < 32 {
					t.Errorf("expected UUID format for Id, got: %s", resp.Id)
				}
				
				if len(resp.Data) != 0 {
					for i, item := range resp.Data {
						if item == "" {
							t.Errorf("Data[%d] is empty", i)
						}
					}
				}
			},
		},
		{
			name:           "Invalid Quality Parameter Test",
			density:        "300",
			quality:        "invalid",
			fileContent:    "Test PDF Content",
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *apiV1.ConvertResponse) {
				if resp.Id == "" {
					t.Errorf("expected non-empty ID, got empty string")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup dependencies for testing
			fileBuilder := builder.NewFileBuilder()
			imageConvertService := service.NewImageMagickConvertService()
			convertUsecase := usecase.NewConvertUsecase(fileBuilder, imageConvertService)
			apiV1Service := v1.NewService(convertUsecase)
			server := app.NewServer(apiV1Service)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			
			_ = writer.WriteField("density", tt.density)
			_ = writer.WriteField("quality", tt.quality)
			
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

			req := httptest.NewRequest("POST", "/v1/convert", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, recorder.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp apiV1.ConvertResponse
				respBody := recorder.Body.Bytes()
				
				if !json.Valid(respBody) {
					t.Fatalf("response is not valid JSON: %s", respBody)
				}
				
				err := json.Unmarshal(respBody, &resp)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				
				if resp.Id == "" {
					t.Errorf("expected non-empty ID, got empty string")
				}
				
				if resp.Data == nil {
					t.Errorf("expected non-nil Data array, got nil")
				}
				
				contentType := recorder.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("expected Content-Type to be application/json, got %s", contentType)
				}
				
				if tt.validateResp != nil {
					tt.validateResp(t, &resp)
				}
			}
		})
	}
}
