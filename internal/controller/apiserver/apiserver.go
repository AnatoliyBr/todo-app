package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errIncorrectAuthHeader      = errors.New("incorrect auth header")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey uint8

type server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	uc     usecase.UseCase
}

func NewServer(config *Config, uc usecase.UseCase) *server {
	s := &server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
		uc:     uc,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {

	// middleware
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	// test
	s.router.HandleFunc("/hello", s.handleHello()).Methods(http.MethodGet)

	// public
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/tokens", s.handleTokensCreate()).Methods(http.MethodPost)

	// private
	s.router.Handle("/profile", s.authenticateUser(s.handleUserProfile())).Methods(http.MethodGet)

	listSubrouter := s.router.PathPrefix("/lists").Subrouter()
	listSubrouter.Use(s.authenticateUser)
	listSubrouter.HandleFunc("", s.handleListsCreate()).Methods(http.MethodPost)
	listSubrouter.HandleFunc("", s.handleListsGetByUser()).Methods(http.MethodGet)
}

func (s *server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *server) StartServer() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.BindAddr, s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]string{"test": "hello"})
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	type tokenClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			s.error(w, r, http.StatusBadRequest, errIncorrectAuthHeader)
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" || authHeaderParts[1] == "" {
			s.error(w, r, http.StatusBadRequest, errIncorrectAuthHeader)
			return
		}

		tokenString := authHeaderParts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&tokenClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("incorrect signing method")
				}
				return []byte(s.config.SecretKey), nil
			})

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		claims := token.Claims.(*tokenClaims)

		u, err := s.uc.UsersFindByID(claims.UserID)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
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

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.config.SecretKey))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w,
			&http.Cookie{
				Name:  "token",
				Value: tokenString,
			})

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser))
	}
}

func (s *server) handleListsCreate() http.HandlerFunc {
	type request struct {
		ListTitle string `json:"list_title"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(ctxKeyUser).(*entity.User)

		l := &entity.List{
			ListTitle: req.ListTitle,
			UserID:    u.UserID,
		}

		if err := s.uc.ListsCreate(l); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, l)
	}
}

func (s *server) handleListsGetByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*entity.User)

		list, err := s.uc.ListsFindByUser(u.UserID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, list)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(data)
	}
}
