package usecase

import (
	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
)

type AppUseCase struct {
	store store.Store
}

func NewAppUseCase(s store.Store) *AppUseCase {
	return &AppUseCase{
		store: s,
	}
}

func (uc *AppUseCase) UsersCreate(u *entity.User) error {
	return uc.store.User().Create(u)
}

func (uc *AppUseCase) UsersFindByID(id int) (*entity.User, error) {
	return uc.store.User().FindByID(id)
}

func (uc *AppUseCase) UsersFindByEmail(email string) (*entity.User, error) {
	return uc.store.User().FindByEmail(email)
}
