package api

import (
	"github.com/mohammadsf7293/blank-go-project/internal/db"
)

type Server struct {
	store db.Store
}

func NewServer(store db.Store) *Server {
	return &Server{
		store: store,
	}
}
