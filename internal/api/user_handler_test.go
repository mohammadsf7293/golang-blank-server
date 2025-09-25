package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mohammadsf7293/blank-go-project/internal/db"
	"github.com/mohammadsf7293/blank-go-project/internal/mocks"
)

func TestHandleUsers_Get(t *testing.T) {
	store := mocks.NewStore()
	// Add test users
	store.CreateUser(context.Background(), db.CreateUserParams{
		Username:     "testuser1",
		Email:        "test1@example.com",
		PasswordHash: "hash1",
	})
	store.CreateUser(context.Background(), db.CreateUserParams{
		Username:     "testuser2",
		Email:        "test2@example.com",
		PasswordHash: "hash2",
	})

	server := NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/api/users?limit=10&offset=0", nil)
	w := httptest.NewRecorder()

	server.HandleUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var users []db.User
	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

func TestHandleUsers_Post(t *testing.T) {
	store := mocks.NewStore()
	server := NewServer(store)

	input := `{"username":"newuser","email":"new@example.com","password":"pass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.HandleUsers(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]int64
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if id, ok := response["id"]; !ok || id != 1 {
		t.Errorf("Expected user ID 1, got %v", id)
	}

	// Verify user was added to store
	users, _ := store.ListUsers(context.Background(), db.ListUsersParams{Limit: 10})
	if len(users) != 1 {
		t.Errorf("Expected 1 user in store, got %d", len(users))
	}
}
