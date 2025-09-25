package db

import (
	"context"
	"database/sql"
)

//go:generate mockgen -destination=../mocks/mock_store.go -package=mocks github.com/mohammadsf7293/blank-go-project/internal/db Store

// Store defines the interface for database operations
type Store interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	Close() error
}
