package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"todo/internal/usecase"
)

type Middleware interface {
	Logger(next http.Handler) http.Handler
	Auth(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	usecaseSession usecase.Session
}

func New(usecaseSession usecase.Session) Middleware {
	return &middleware{usecaseSession: usecaseSession}
}

func (m *middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, r.UserAgent(), r.RemoteAddr)
	})
}

func (m *middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Println("error get cookie", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println(cookie)
		session, err := m.usecaseSession.Get(r.Context(), cookie.Value)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || errors.Is(err, usecase.ErrSessionHasExpired) {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next(w, r.WithContext(ctx))
	}
}
