package testrepository

import (
	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
)

type UserRepository struct {
	users map[int]*entity.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int]*entity.User),
	}
}

func (r *UserRepository) Create(u *entity.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.UserID = len(r.users) + 1
	r.users[u.UserID] = u

	return nil
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}
