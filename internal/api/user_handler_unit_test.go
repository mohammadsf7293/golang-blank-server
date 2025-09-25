package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mohammadsf7293/blank-go-project/internal/db"
	"github.com/mohammadsf7293/blank-go-project/internal/mocks"
	"go.uber.org/mock/gomock"
)

type mockResult struct {
	lastID int64
}

func (m mockResult) LastInsertId() (int64, error) {
	return m.lastID, nil
}

func (m mockResult) RowsAffected() (int64, error) {
	return 1, nil
}

func TestHandleUsers_GetUnit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockStore)
		query         string
		expectedCode   int
		expectedUsers  int
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockStore) {
				m.EXPECT().
					ListUsers(gomock.Any(), db.ListUsersParams{Limit: 10, Offset: 0}).
					Return([]db.User{
						{ID: 1, Username: "user1", Email: "user1@example.com"},
						{ID: 2, Username: "user2", Email: "user2@example.com"},
					}, nil)
			},
			query: "?limit=10&offset=0",
			expectedCode: http.StatusOK,
			expectedUsers: 2,
		},
		{
			name: "database error",
			setupMock: func(m *mocks.MockStore) {
				m.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			query: "?limit=10&offset=0",
			expectedCode: http.StatusInternalServerError,
			expectedUsers: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewMockStore(ctrl)
			tt.setupMock(mockStore)

			server := NewServer(mockStore)
			req := httptest.NewRequest(http.MethodGet, "/api/users"+tt.query, nil)
			w := httptest.NewRecorder()

			server.HandleUsers(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				var users []db.User
				if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if len(users) != tt.expectedUsers {
					t.Errorf("Expected %d users, got %d", tt.expectedUsers, len(users))
				}
			}
		})
	}
}

func TestHandleUsers_PostUnit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockStore)
		input          string
		expectedCode   int
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockStore) {
				m.EXPECT().
					CreateUser(gomock.Any(), db.CreateUserParams{
						Username: "testuser",
						Email:    "test@example.com",
						PasswordHash: "password123",
					}).
					Return(mockResult{lastID: 1}, nil)
			},
			input: `{"username":"testuser","email":"test@example.com","password":"password123"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "database error",
			setupMock: func(m *mocks.MockStore) {
				m.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			input: `{"username":"testuser","email":"test@example.com","password":"password123"}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "invalid json",
			setupMock: func(m *mocks.MockStore) {},
			input: `{invalid json}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewMockStore(ctrl)
			tt.setupMock(mockStore)

			server := NewServer(mockStore)
			req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.HandleUsers(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestHandleUsers_MethodNotAllowedUnit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockStore(ctrl)
	server := NewServer(mockStore)

	methods := []string{http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/users", nil)
			w := httptest.NewRecorder()

			server.HandleUsers(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status code %d for method %s, got %d",
					http.StatusMethodNotAllowed, method, w.Code)
			}
		})
	}
}