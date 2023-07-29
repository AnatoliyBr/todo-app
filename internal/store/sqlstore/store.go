package sqlstore

import (
	"database/sql"

	"github.com/AnatoliyBr/todo-app/internal/store"
)

type Store struct {
	UserRepository *UserRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		UserRepository: NewUserRepository(db),
	}
}

func (s *Store) User() store.UserRepository {
	return s.UserRepository
}
