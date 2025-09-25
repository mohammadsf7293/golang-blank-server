package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mohammadsf7293/blank-go-project/internal/api"
	"github.com/mohammadsf7293/blank-go-project/internal/config"
	"github.com/mohammadsf7293/blank-go-project/internal/db"
)

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
	srv := api.NewServer(store)

	// Initialize router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/health", srv.HandleHealth)
	mux.HandleFunc("/api/users", srv.HandleUsers)

	// Start server
	log.Printf("Starting server on :%s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}
