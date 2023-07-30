package entity_test

import (
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *entity.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *entity.User {
				return entity.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "encryptedpassword"
				return u
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *entity.User {
				u := entity.TestUser(t)
				u.Password = "short"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := entity.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}
