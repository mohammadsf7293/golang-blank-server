package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mohammadsf7293/blank-go-project/internal/config"
	"github.com/mohammadsf7293/blank-go-project/internal/db"
)

type server struct {
	store *db.Store
}

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize database
	store, err := db.NewStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// Initialize server
	srv := &server{store: store}

	// Initialize router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/health", srv.handleHealth)
	mux.HandleFunc("/api/users", srv.handleUsers)

	// Start server
	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (s *server) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// List users
		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
		offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 32)
		if limit == 0 {
			limit = 10
		}

		users, err := s.store.ListUsers(r.Context(), db.ListUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		// Create user
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := s.store.CreateUser(r.Context(), db.CreateUserParams{
			Username:     input.Username,
			Email:       input.Email,
			PasswordHash: input.Password, // Note: In production, hash the password
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{"id": id})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}