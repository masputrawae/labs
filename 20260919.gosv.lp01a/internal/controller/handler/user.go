package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"todo/internal/domain/model"
	"todo/internal/usecase"

	"github.com/go-playground/validator/v10"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type User interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type user struct {
	usecaseUser    usecase.User
	usecaseSession usecase.Session
	validate       *validator.Validate
}

func NewUser(
	usecaseUser usecase.User,
	usecaseSession usecase.Session,
	validate *validator.Validate,
) User {
	return &user{
		usecaseUser:    usecaseUser,
		usecaseSession: usecaseSession,
		validate:       validate,
	}
}

func (u *user) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if !Bind(w, r, u.validate, &req) {
		return
	}

	id, err := u.usecaseUser.Register(r.Context(), req)
	if err != nil {
		sqlErr, ok := errors.AsType[*sqlite.Error](err)
		if ok {
			switch sqlErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				s := err.Error()
				switch {
				case strings.Contains(s, "users.username"):
					JSON(w, http.StatusConflict, H{"error": H{"username": "CONFLICT"}})
				case strings.Contains(s, "users.email"):
					JSON(w, http.StatusConflict, H{"error": H{"email": "CONFLICT"}})
				}
				return
			}
		}

		log.Println("error user register:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := u.usecaseSession.Create(r.Context(), id)
	if err != nil {
		log.Println("error create session:", err)
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (u *user) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if !Bind(w, r, u.validate, &req) {
		return
	}

	id, err := u.usecaseUser.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println("error user login:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	cookie, err := u.usecaseSession.Create(r.Context(), id)
	if err != nil {
		log.Println("error create session:", err)
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (u *user) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return
	}

	sessionID := cookie.Value
	if err = u.usecaseSession.Delete(r.Context(), sessionID); err != nil {
		log.Println("error delete session:", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-1).UTC(),
	})

	w.WriteHeader(http.StatusNoContent)
}

func (u *user) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	req := model.UserUpdateRequest{
		UserID: userID,
	}

	if !Bind(w, r, u.validate, &req) {
		return
	}

	if err := u.usecaseUser.Update(r.Context(), req); err != nil {
		if errors.Is(err, usecase.ErrCredentials) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		sqlErr, ok := errors.AsType[*sqlite.Error](err)
		if ok {
			switch sqlErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				s := err.Error()
				switch {
				case strings.Contains(s, "users.username"):
					JSON(w, http.StatusConflict, H{"error": H{"username": "CONFLICT"}})
				case strings.Contains(s, "users.email"):
					JSON(w, http.StatusConflict, H{"error": H{"email": "CONFLICT"}})
				}
				return
			}

		}
		log.Println("error user update:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *user) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	data, err := u.usecaseUser.Get(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("error get user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": data})
}

func (u *user) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	req := model.UserDeleteRequest{
		UserID: userID,
	}
	if !Bind(w, r, u.validate, &req) {
		return
	}

	if err := u.usecaseUser.Delete(r.Context(), req); err != nil {
		if errors.Is(err, usecase.ErrCredentials) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		MaxAge:  -1,
		Expires: time.Now().Add(-1).UTC(),
	})
	w.WriteHeader(http.StatusNoContent)
}
