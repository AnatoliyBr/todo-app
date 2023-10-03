package usecase

import "github.com/AnatoliyBr/todo-app/internal/entity"

type UseCase interface {
	UsersCreate(*entity.User) error
	UsersFindByID(int) (*entity.User, error)
	UsersFindByEmail(string) (*entity.User, error)

	ListsCreate(*entity.List) error
	ListsFindByID(int, int) (*entity.List, error)
	ListsEdit(*entity.List) (*entity.List, error)
	ListsDelete(*entity.List) error
	ListsFindByUser(int) ([]*entity.List, error)
}
