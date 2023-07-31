package sqlstore

import (
	"database/sql"

	"github.com/AnatoliyBr/todo-app/internal/store"
)

type Store struct {
	userRepository *UserRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		userRepository: NewUserRepository(db),
	}
}

func (s *Store) User() store.UserRepository {
	return s.userRepository
}
