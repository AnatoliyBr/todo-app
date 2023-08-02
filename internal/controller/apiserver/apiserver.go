package apiserver

import (
	"io"
	"log"
	"net/http"

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
