package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
	"github.com/AnatoliyBr/todo-app/internal/store/testrepository"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleHello(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, rec.Body)
}

func TestServer_SetRequestID(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	s.setRequestID(handler).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
}
func TestServer_AuthenticateUser(t *testing.T) {
	type tokenClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}

	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	u := entity.TestUser(t)
	s.uc.UsersCreate(u)

	testCases := []struct {
		name         string
		tokenString  func() string
		expectedCode int
	}{
		{
			name: "authenticated",
			tokenString: func() string {
				claims := &tokenClaims{
					UserID: u.UserID,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(s.config.SecretKey))
				return tokenString
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "incorrect auth header",
			tokenString: func() string {
				return ""
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not authenticated",
			tokenString: func() string {
				claims := &tokenClaims{}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(s.config.SecretKey))
				return tokenString
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "expired token",
			tokenString: func() string {
				claims := &tokenClaims{
					UserID: u.UserID,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute * 5)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(s.config.SecretKey))
				return tokenString
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.tokenString()))

			s.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
func TestServer_HandleUsersCreate(t *testing.T) {
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
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

func TestServer_HandleTokensCreate(t *testing.T) {
	u := entity.TestUser(t)
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	s.uc.UsersCreate(u)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/tokens", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUserProfile(t *testing.T) {
	u := entity.TestUser(t)
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	s.uc.UsersCreate(u)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	req.WithContext(context.WithValue(req.Context(), ctxKeyUser, u))

	s.handleUserProfile().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotNil(t, rec.Body)
}
