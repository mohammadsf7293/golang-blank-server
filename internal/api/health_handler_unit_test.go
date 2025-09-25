package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mohammadsf7293/blank-go-project/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestHandleHealth_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		method       string
		expectedCode int
	}{
		{
			name:         "success - GET method",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			name:         "method not allowed - POST",
			method:       http.MethodPost,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "method not allowed - PUT",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "method not allowed - DELETE",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewMockStore(ctrl)
			server := NewServer(mockStore)

			req := httptest.NewRequest(tt.method, "/api/health", nil)
			w := httptest.NewRecorder()

			server.HandleHealth(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if status, ok := response["status"]; !ok || status != "healthy" {
					t.Errorf("Expected status 'healthy', got %v", status)
				}
			}
		})
	}
}