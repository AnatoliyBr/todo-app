package sqlrepository_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
	"github.com/AnatoliyBr/todo-app/internal/store/sqlrepository"
	"github.com/stretchr/testify/assert"
)

func TestListRepository_Create(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("users", "lists")

	ur := sqlrepository.NewUserRepository(db)
	lr := sqlrepository.NewListRepository(db)
	s := store.NewAppStore(ur, lr)
	u := entity.TestUser(t)
	l := entity.TestList(t)
	l.UserID = 10

	err := s.List().Create(l)
	assert.Error(t, err)

	s.User().Create(u)
	l.UserID = u.UserID
	err = s.List().Create(l)
	assert.NoError(t, err)
}

func TestListRepository_FindByTitle(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("users", "lists")

	ur := sqlrepository.NewUserRepository(db)
	lr := sqlrepository.NewListRepository(db)
	s := store.NewAppStore(ur, lr)
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	s.User().Create(u)
	l1.UserID = u.UserID

	_, err := s.List().FindByTitle(u.UserID, l1.ListTitle)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.List().Create(l1)
	l2, err := s.List().FindByTitle(u.UserID, l1.ListTitle)
	assert.NoError(t, err)
	assert.NotNil(t, l2)
}

func TestListRepository_Edit(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("users", "lists")

	ur := sqlrepository.NewUserRepository(db)
	lr := sqlrepository.NewListRepository(db)
	s := store.NewAppStore(ur, lr)
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	s.User().Create(u)
	l1.UserID = u.UserID
	s.List().Create(l1)

	l1.ListTitle = "TEST TITLE 2"
	l2, err := s.List().Edit(l1)
	assert.NoError(t, err)
	assert.NotNil(t, l2)
}

func TestListRepository_Delete(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("users", "lists")

	ur := sqlrepository.NewUserRepository(db)
	lr := sqlrepository.NewListRepository(db)
	s := store.NewAppStore(ur, lr)
	u := entity.TestUser(t)
	l := entity.TestList(t)
	s.User().Create(u)
	l.UserID = u.UserID
	s.List().Create(l)

	err := s.List().Delete(l)
	assert.NoError(t, err)

	_, err = s.List().FindByTitle(u.UserID, l.ListTitle)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}

func TestListRepository_FindByUser(t *testing.T) {
	db, teardown := sqlrepository.TestDB(t, testDatabaseURL)
	defer teardown("users", "lists")

	ur := sqlrepository.NewUserRepository(db)
	lr := sqlrepository.NewListRepository(db)
	s := store.NewAppStore(ur, lr)
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	l2 := entity.TestList(t)
	l2.ListTitle = "TEST TITLE 2"
	s.User().Create(u)
	l1.UserID = u.UserID
	l2.UserID = u.UserID

	_, err := s.List().FindByUser(u.UserID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.List().Create(l1)
	s.List().Create(l2)
	lists, err := s.List().FindByUser(u.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, lists)
}
