package db

import (
	"context"
	"database/sql"
)

// Store defines the interface for database operations
type Store interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	Close() error
}
