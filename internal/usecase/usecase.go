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

func (uc *AppUseCase) ListsCreate(l *entity.List) error {
	return uc.store.List().Create(l)
}

func (uc *AppUseCase) ListsFindByID(listID, userID int) (*entity.List, error) {
	return uc.store.List().FindByID(listID, userID)
}

func (uc *AppUseCase) ListsEdit(l *entity.List) (*entity.List, error) {
	return uc.store.List().Edit(l)
}

func (uc *AppUseCase) ListsDelete(l *entity.List) error {
	return uc.store.List().Delete(l)
}

func (uc *AppUseCase) ListsFindByUser(userID int) ([]*entity.List, error) {
	return uc.store.List().FindByUser(userID)
}
