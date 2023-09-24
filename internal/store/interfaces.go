package store

import "github.com/AnatoliyBr/todo-app/internal/entity"

type Store interface {
	User() UserRepository
	List() ListRepository
}

type UserRepository interface {
	Create(*entity.User) error
	FindByID(int) (*entity.User, error)
	FindByEmail(string) (*entity.User, error)
}

type ListRepository interface {
	Create(*entity.List) error
	FindByTitle(int, string) (*entity.List, error)
	Edit(*entity.List) (*entity.List, error)
	Delete(*entity.List) error
	FindByUser(int) ([]*entity.List, error)
}
