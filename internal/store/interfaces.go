package store

import "github.com/AnatoliyBr/todo-app/internal/entity"

type Store interface {
	User() UserRepository
}

type UserRepository interface {
	Create(*entity.User) error
	FindByID(int) (*entity.User, error)
	FindByEmail(string) (*entity.User, error)
}
