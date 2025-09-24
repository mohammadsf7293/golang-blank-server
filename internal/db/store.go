package db

import (
	"database/sql"

	"github.com/mohammadsf7293/blank-go-project/internal/config"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(cfg *config.Config) (*Store, error) {
	db, err := sql.Open("mysql", cfg.DBConfig.DSN())
	if err != nil {
		return nil, err
	}

	return &Store{
		Queries: New(db),
		db:      db,
	}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}