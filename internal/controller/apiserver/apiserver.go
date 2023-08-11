package apiserver

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

type server struct {
	config *Config
	router *mux.Router
	uc     usecase.UseCase
}

func NewServer(config *Config, uc usecase.UseCase) *server {
	s := &server{
		config: config,
		router: mux.NewRouter(),
		uc:     uc,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello()).Methods(http.MethodGet)

	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/tokens", s.handleTokensCreate()).Methods(http.MethodPost)
}

func (s *server) StartServer() error {
	log.Println("start server...")
	return http.ListenAndServe(s.config.BindAddr, s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &entity.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.uc.UsersCreate(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleTokensCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type tokenClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.uc.UsersFindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		exp := time.Now().Add(time.Minute * 5)
		claims := &tokenClaims{
			UserID: u.UserID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		}

		t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		st, err := t.SignedString([]byte(s.config.SecretKey))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w,
			&http.Cookie{
				Name:  "token",
				Value: st,
			})

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
