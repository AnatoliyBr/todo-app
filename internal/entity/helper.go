package entity

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestList(t *testing.T) *List {
	return &List{
		ListTitle: "TEST TITLE 1",
	}
}
