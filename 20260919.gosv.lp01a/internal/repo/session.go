package repo

import (
	"context"
	"database/sql"
	"time"
	"todo/internal/domain/entity"
)

type Session interface {
	Create(ctx context.Context, arg SessionCreateParams) error
	Get(ctx context.Context, sessionID string) (*entity.Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteHasExpired(ctx context.Context, now time.Time) error
}

type session struct {
	*sql.DB
}

func NewSession(db *sql.DB) Session {
	return &session{db}
}

type SessionCreateParams struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

func (s *session) Create(ctx context.Context, arg SessionCreateParams) error {
	q := `
	INSERT INTO sessions
		( id, expires_at, user_id)
	VALUES ( ?, ?, ? )
	`
	_, err := s.ExecContext(ctx, q, arg.ID, arg.ExpiresAt, arg.UserID)
	return err
}

func (s *session) Get(ctx context.Context, sessionID string) (*entity.Session, error) {
	q := "SELECT id, expires_at, user_id FROM sessions WHERE id = ?"

	i := new(entity.Session)
	if err := s.
		QueryRowContext(ctx, q, sessionID).
		Scan(
			&i.ID,
			&i.ExpiresAt,
			&i.UserID,
		); err != nil {
		return nil, err
	}

	return i, nil
}

func (s *session) Delete(ctx context.Context, sessionID string) error {
	q := "DELETE FROM sessions WHERE id = ?"
	_, err := s.ExecContext(ctx, q, sessionID)
	return err
}

func (s *session) DeleteHasExpired(ctx context.Context, now time.Time) error {
	q := "DELETE FROM sessions WHERE expires_at <= ?"
	_, err := s.ExecContext(ctx, q, now)
	return err
}
