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
	"github.com/gorilla/mux"
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
	u1 := entity.TestUser(t)
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	s.uc.UsersCreate(u1)
	u1.Sanitize()
	u1.EncryptedPassword = ""

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, u1))

	s.handleUserProfile().ServeHTTP(rec, req)

	u2 := &entity.User{}
	json.NewDecoder(rec.Body).Decode(u2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, u1, u2)
}

func TestSerer_HandleListsCreate(t *testing.T) {
	u := entity.TestUser(t)
	l := entity.TestList(t)
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
				"list_title": l.ListTitle,
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid title",
			payload: map[string]string{
				"list_title": "inv@li_d title *",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/lists", b)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, u))

			s.handleListsCreate().ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleListsGetByUser(t *testing.T) {
	u := entity.TestUser(t)
	l1 := entity.TestList(t)
	l2 := entity.TestList(t)
	l2.ListTitle = "TEST TITLE 2"
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	s.uc.UsersCreate(u)
	l1.UserID = u.UserID
	l2.UserID = u.UserID
	s.uc.ListsCreate(l1)
	s.uc.ListsCreate(l2)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/lists", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, u))

	s.handleListsGetByUser().ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_HandleListsGetByID(t *testing.T) {
	u := entity.TestUser(t)
	l := entity.TestList(t)
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	s.uc.UsersCreate(u)
	l.UserID = u.UserID
	s.uc.ListsCreate(l)

	assert.NotNil(t, l.ListID)

	testCases := []struct {
		name         string
		id           string
		expectedCode int
	}{
		{
			name:         "valid",
			id:           "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "not found",
			id:           "2",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid",
			id:           "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/lists/%s", tc.id), nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, u))
			req = mux.SetURLVars(req, map[string]string{"listID": tc.id})

			s.handleListsGetByID().ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleListsEdit(t *testing.T) {
	u := entity.TestUser(t)
	ur := testrepository.NewUserRepository()
	lr := testrepository.NewListRepository()
	store := store.NewAppStore(ur, lr)
	uc := usecase.NewAppUseCase(store)
	s := NewServer(NewConfig(), uc)
	s.uc.UsersCreate(u)

	testCases := []struct {
		name         string
		id           string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			id:   "1",
			payload: map[string]string{
				"list_title": "NEW TITLE 2",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			id:           "1",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			id:   "2",
			payload: map[string]string{
				"list_title": "TEST TITLE 2",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid id",
			id:   "invalid",
			payload: map[string]string{
				"list_title": "TEST TITLE 2",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid title",
			id:   "1",
			payload: map[string]string{
				"list_title": "inv@li_d title *",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "existing title",
			id:   "1",
			payload: map[string]string{
				"list_title": "TEST TITLE 1",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := entity.TestList(t)
			l.UserID = u.UserID
			s.uc.ListsCreate(l)
			assert.NotNil(t, l.ListID)

			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/lists/%s", tc.id), b)
			req = req.WithContext(context.WithValue(req.Context(), ctxKeyUser, u))
			req = mux.SetURLVars(req, map[string]string{"listID": tc.id})

			s.handleListsEdit().ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

			s.uc.ListsDelete(l)
		})
	}
}
