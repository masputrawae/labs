package repo

import (
	"context"
	"database/sql"
	"time"
	"todo/internal/domain/entity"
)

type User interface {
	Create(ctx context.Context, arg UserCreateParams) error
	Update(ctx context.Context, arg UserUpdateParams) error
	Delete(ctx context.Context, userID string) error
	Get(ctx context.Context, userID string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}

type user struct {
	*sql.DB
}

func NewUser(db *sql.DB) User {
	return &user{db}
}

type UserCreateParams struct {
	ID        string
	Email     string
	Name      string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *user) Create(ctx context.Context, arg UserCreateParams) error {
	q := `
	INSERT INTO users
		( id, email, name, username, password, created_at, updated_at )
	VALUES ( ?, ?, ?, ?, ?, ?, ? )`
	_, err := u.ExecContext(ctx, q,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.Username,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

type UserUpdateParams struct {
	ID        string
	Email     *string
	Name      *string
	Username  *string
	Password  *string
	UpdatedAt time.Time
}

func (u *user) Update(ctx context.Context, arg UserUpdateParams) error {
	q := `
	UPDATE users
	SET
		email 			= COALESCE(?1, email),
		name 				= COALESCE(?2, name),
		username 		= COALESCE(?3, username),
		password 		= COALESCE(?4, password),
		updated_at 	= ?5
	WHERE id = ?6
	`
	result, err := u.ExecContext(ctx, q,
		arg.Email,
		arg.Name,
		arg.Username,
		arg.Password,
		arg.UpdatedAt,
		arg.ID,
	)

	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *user) Delete(ctx context.Context, userID string) error {
	q := `DELETE FROM users WHERE id = ?`
	result, err := u.ExecContext(ctx, q, userID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *user) Get(ctx context.Context, userID string) (*entity.User, error) {
	q := `
	SELECT 
		id, name, email, username, password, created_at, updated_at 
	FROM users 
	WHERE id = ?`

	i := new(entity.User)
	if err := u.
		QueryRowContext(ctx, q, userID).
		Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Username,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return i, nil
}

func (u *user) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	q := `
	SELECT 
		id, name, email, username, password, created_at, updated_at 
	FROM users 
	WHERE username = ?`

	i := new(entity.User)
	if err := u.
		QueryRowContext(ctx, q, username).
		Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Username,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return i, nil
}
