package usecase_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
	"github.com/AnatoliyBr/todo-app/internal/store/testrepository"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAppUseCase_UsersCreate(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	u := entity.TestUser(t)
	uc := usecase.NewAppUseCase(s)

	assert.NotNil(t, u)
	assert.NoError(t, uc.UsersCreate(u))
}

func TestAppUseCase_UsersFindByID(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
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
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(s)
	u1 := entity.TestUser(t)
	_, err := uc.UsersFindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	uc.UsersCreate(u1)
	u2, err := uc.UsersFindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestAppUseCase_ListsCreate(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(s)
	u := entity.TestUser(t)
	l := entity.TestList(t)
	uc.UsersCreate(u)
	l.UserID = u.UserID

	err := uc.ListsCreate(l)
	assert.NoError(t, err)
}

func TestAppUseCase_ListsFindByID(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(s)
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	uc.UsersCreate(u)
	l1.UserID = u.UserID

	_, err := uc.ListsFindByID(l1.ListID, u.UserID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	uc.ListsCreate(l1)
	l2, err := uc.ListsFindByID(l1.ListID, u.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, l2)
}

func TestAppUseCase_ListsEdit(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(s)
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	l2 := entity.TestList(t)
	uc.UsersCreate(u)
	l1.UserID = u.UserID
	l2.UserID = u.UserID
	uc.ListsCreate(l1)

	l2.ListTitle = "TEST TITLE 2"
	l3, err := uc.ListsEdit(l2)
	assert.NoError(t, err)
	assert.NotNil(t, l3)
}

func TestAppUseCase_ListsDelete(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(s)
	u := entity.TestUser(t)
	l := entity.TestList(t)
	uc.UsersCreate(u)
	l.UserID = u.UserID
	uc.ListsCreate(l)

	err := uc.ListsDelete(l)
	assert.NoError(t, err)

	_, err = uc.ListsFindByID(l.ListID, u.UserID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}

func TestAppUseCase_ListsFindByUser(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	s := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(s)
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	l2 := entity.TestList(t)
	l2.ListTitle = "TEST TITLE 2"
	uc.UsersCreate(u)
	l1.UserID = u.UserID
	l2.UserID = u.UserID

	_, err := uc.ListsFindByUser(u.UserID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	uc.ListsCreate(l1)
	uc.ListsCreate(l2)
	lists, err := uc.ListsFindByUser(u.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, lists)
}
