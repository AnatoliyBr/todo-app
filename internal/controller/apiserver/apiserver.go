package apiserver

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	config *Config
	router *mux.Router
}

func NewServer(config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		config: config,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello()).Methods(http.MethodGet)
}

func (s *server) StartServer() error {
	log.Println("start server...")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	}
}
