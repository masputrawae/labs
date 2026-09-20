package usecase

import (
	"context"
	"errors"
	"net/http"
	"time"
	"todo/internal/domain/entity"
	"todo/internal/repo"
	"uuid"
)

var (
	ErrSessionHasExpired = errors.New("session has expired")
)

type Session interface {
	Create(ctx context.Context, userID string) (*http.Cookie, error)
	Get(ctx context.Context, sessionID string) (*entity.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteHasExpired(ctx context.Context) error
}

type session struct {
	repoSession repo.Session
}

func NewSession(repoSession repo.Session) Session {
	return &session{repoSession: repoSession}
}

func (s *session) Create(ctx context.Context, userID string) (*http.Cookie, error) {
	sessionID := uuid.NewV7().String()
	ttl := 30 * time.Minute
	expiresAt := time.Now().Add(ttl).UTC()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(ttl.Seconds()),
		Secure:   false, //dev only
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}

	err := s.repoSession.Create(ctx, repo.SessionCreateParams{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})

	return cookie, err
}

func (s *session) Get(ctx context.Context, sessionID string) (*entity.Session, error) {
	i, err := s.repoSession.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if i.ExpiresAt.Before(time.Now().UTC()) {
		if err = s.repoSession.Delete(ctx, sessionID); err != nil {
			return nil, err
		}

		return nil, ErrSessionHasExpired
	}

	return i, nil
}

func (s *session) Delete(ctx context.Context, sessionID string) error {
	return s.repoSession.Delete(ctx, sessionID)
}

func (s *session) DeleteHasExpired(ctx context.Context) error {
	return s.repoSession.DeleteHasExpired(ctx, time.Now().UTC())
}
