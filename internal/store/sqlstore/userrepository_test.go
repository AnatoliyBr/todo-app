package sqlstore_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
	"github.com/AnatoliyBr/todo-app/internal/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, testDatabaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	u := entity.TestUser(t)

	assert.NotNil(t, u)
	assert.NoError(t, s.User().Create(u))
}

func TestUserRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, testDatabaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	u1 := entity.TestUser(t)
	_, err := s.User().FindByID(u1.UserID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(u1)
	u2, err := s.User().FindByID(u1.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, testDatabaseURL)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	u1 := entity.TestUser(t)
	_, err := s.User().FindByEmail(u1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(u1)
	u2, err := s.User().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
