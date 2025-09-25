package mocks

import (
	"context"
	"database/sql"
	"time"

	"github.com/mohammadsf7293/blank-go-project/internal/db"
)

type MockResult struct {
	LastID int64
}

func (m MockResult) LastInsertId() (int64, error) {
	return m.LastID, nil
}

func (m MockResult) RowsAffected() (int64, error) {
	return 1, nil
}

// Store implements db.Store interface for testing
type Store struct {
	users []db.User
}

var _ db.Store = (*Store)(nil) // Verify Store implements db.Store

func NewStore() *Store {
	return &Store{
		users: make([]db.User, 0),
	}
}

func (m *Store) Close() error {
	return nil
}

func (m *Store) CreateUser(ctx context.Context, arg db.CreateUserParams) (sql.Result, error) {
	user := db.User{
		ID:           int64(len(m.users) + 1),
		Username:     arg.Username,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	m.users = append(m.users, user)
	return MockResult{LastID: user.ID}, nil
}

func (m *Store) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	start := int(arg.Offset)
	end := start + int(arg.Limit)
	if end > len(m.users) {
		end = len(m.users)
	}
	if start > end {
		return []db.User{}, nil
	}
	return m.users[start:end], nil
}
