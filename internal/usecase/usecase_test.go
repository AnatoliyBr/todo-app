package usecase_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
	"github.com/AnatoliyBr/todo-app/internal/store/teststore"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAppUseCase_UsersCreate(t *testing.T) {
	s := teststore.NewStore()
	u := entity.TestUser(t)
	uc := usecase.NewAppUseCase(s)

	assert.NotNil(t, u)
	assert.NoError(t, uc.UsersCreate(u))
}

func TestAppUseCase_UsersFindByID(t *testing.T) {
	s := teststore.NewStore()
	uc := usecase.NewAppUseCase(s)
	u1 := entity.TestUser(t)
	_, err := uc.UsersFindByID(u1.UserID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	uc.UsersCreate(u1)
	u2, err := uc.UsersFindByID(u1.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestAppUseCase_UsersFindByEmail(t *testing.T) {
	s := teststore.NewStore()
	uc := usecase.NewAppUseCase(s)
	u1 := entity.TestUser(t)
	_, err := uc.UsersFindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	uc.UsersCreate(u1)
	u2, err := uc.UsersFindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
