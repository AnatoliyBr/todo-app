package sqlstore_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
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
