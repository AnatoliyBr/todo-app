package apiserver

import (
	"bytes"
	"encoding/json"
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

func TestServer_HandleUsersCreate(t *testing.T) {
	uc := usecase.NewAppUseCase(teststore.NewStore())
	s := NewServer(NewConfig(), uc)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
