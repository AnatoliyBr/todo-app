package sqlrepository_test

import (
	"os"
	"testing"
)

var (
	testDatabaseURL string
)

func TestMain(m *testing.M) {
	testDatabaseURL = os.Getenv("TEST_DATABSE_URL_LOCALHOST")
	if testDatabaseURL == "" {
		testDatabaseURL = "postgres://dev:qwerty@localhost:5432/todo_test"
	}

	os.Exit(m.Run())
}
