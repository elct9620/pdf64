package main

import (
	"bytes"
	"context"
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
	"github.com/elct9620/pdf64/internal/entity"
	"github.com/elct9620/pdf64/internal/usecase"
	apiV1 "github.com/elct9620/pdf64/pkg/apis/v1"
)

// MockImageConvertService is a mock implementation of the ImageConvertService interface
type MockImageConvertService struct{}

// Convert always returns a successful result with a mock base64 image
func (m *MockImageConvertService) Convert(ctx context.Context, file *entity.File, options usecase.ImageConvertOptions) ([]string, error) {
	// Return a mock base64 encoded image
	return []string{"data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAHwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIhMUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXGx8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAVYnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOEhYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq8vP09fb3+Pn6/9oADAMBAAIRAxEAPwD3+iiigD//2Q=="}, nil
}

// MockPdfDecryptService is a mock implementation of the PdfDecryptService interface
type MockPdfDecryptService struct{}

// Decrypt always returns success
func (m *MockPdfDecryptService) Decrypt(ctx context.Context, file *entity.File, password string) error {
	// Mark the file as decrypted
	file.Decrypt()
	return nil
}

func TestApiV1Convert(t *testing.T) {
	tests := []struct {
		name           string
		density        string
		quality        string
		password       string
		fileContent    string
		expectedStatus int
		validateResp   func(t *testing.T, resp *apiV1.ConvertResponse)
	}{
		{
			name:           "Basic Conversion Test",
			density:        "300",
			quality:        "90",
			password:       "",
			fileContent:    "%PDF-1.5\n%%EOF\n", // Minimal valid PDF structure
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
			name:           "Password Protected PDF Test",
			density:        "300",
			quality:        "90",
			password:       "secret123",
			fileContent:    "%PDF-1.5\n%%EOF\n", // Minimal valid PDF structure
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, resp *apiV1.ConvertResponse) {
				if len(resp.Id) < 32 {
					t.Errorf("expected UUID format for Id, got: %s", resp.Id)
				}
				
				if len(resp.Data) == 0 {
					t.Errorf("expected non-empty Data array")
				}
			},
		},
		{
			name:           "Invalid Quality Parameter Test",
			density:        "300",
			quality:        "invalid",
			password:       "",
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
			// Setup dependencies for testing with mock services
			fileBuilder := builder.NewFileBuilder()
			
			// Create mock services
			mockImageConvertService := &MockImageConvertService{}
			mockPdfDecryptService := &MockPdfDecryptService{}
			
			convertUsecase := usecase.NewConvertUsecase(fileBuilder, mockImageConvertService, mockPdfDecryptService)
			apiV1Service := v1.NewService(convertUsecase)
			server := app.NewServer(apiV1Service)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			
			_ = writer.WriteField("density", tt.density)
			_ = writer.WriteField("quality", tt.quality)
			if tt.password != "" {
				_ = writer.WriteField("password", tt.password)
			}
			
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
