package apiserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/gorilla/mux"
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

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
