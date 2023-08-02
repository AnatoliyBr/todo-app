package usecase

import "github.com/AnatoliyBr/todo-app/internal/entity"

type UseCase interface {
	UsersCreate(*entity.User) error
	UsersFindByID(int) (*entity.User, error)
	UsersFindByEmail(string) (*entity.User, error)
}
