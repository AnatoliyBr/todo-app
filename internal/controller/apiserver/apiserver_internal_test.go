package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatoliyBr/todo-app/internal/store/teststore"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleHello(t *testing.T) {
	uc := usecase.NewAppUseCase(teststore.NewStore())
	s := NewServer(NewConfig(), uc)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	// s.handleHello().ServeHTTP(rec, req) // #1
	s.ServeHTTP(rec, req) // #2
	assert.Equal(t, "hello", rec.Body.String())
}
