package teststore

import "github.com/AnatoliyBr/todo-app/internal/store"

type Store struct {
	userRepository *UserRepository
}

func NewStore() *Store {
	return &Store{
		userRepository: NewUserRepository(),
	}
}

func (s *Store) User() store.UserRepository {
	return s.userRepository
}
