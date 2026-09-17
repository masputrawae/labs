package repo

import (
	"context"
	"database/sql"
	"time"
	"todo/internal/model"
)

type Repo interface {
	Create(ctx context.Context, task string) (*model.Todo, error)
	Update(ctx context.Context, todoID int64, arg model.TodoUpdateRequest) (*model.Todo, error)
	GetAll(ctx context.Context) ([]model.Todo, error)
	Delete(ctx context.Context, todoID int64) error
	DeleteAll(ctx context.Context) error
}

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, task string) (*model.Todo, error) {
	q := `
	INSERT INTO todos ( task, is_done, created_at, updated_at )
	VALUES ( ?, ?, ?, ? )
	RETURNING id, task, is_done, created_at, updated_at
	`

	now := time.Now().UTC()
	i := new(model.Todo)
	err := r.db.QueryRowContext(ctx, q, task, false, now, now).
		Scan(
			&i.ID,
			&i.Task,
			&i.IsDone,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
	return i, err
}

func (r *repo) Update(ctx context.Context, todoID int64, arg model.TodoUpdateRequest) (*model.Todo, error) {
	q := `
	UPDATE todos
	SET
		task 				= COALESCE(?1, task),
		is_done 		= COALESCE(?2, is_done),
		updated_at 	= ?3
	WHERE id = ?4
	RETURNING id, task, is_done, created_at, updated_at
	`

	now := time.Now().UTC()
	i := new(model.Todo)
	err := r.db.QueryRowContext(ctx, q, arg.Task, arg.IsDone, now, todoID).
		Scan(
			&i.ID,
			&i.Task,
			&i.IsDone,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
	return i, err
}

func (r *repo) GetAll(ctx context.Context) ([]model.Todo, error) {
	q := `SELECT id, task, is_done, created_at, updated_at FROM todos`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items = []model.Todo{}
	for rows.Next() {
		var i model.Todo
		err = rows.Scan(
			&i.ID,
			&i.Task,
			&i.IsDone,
			&i.CreatedAt,
			&i.UpdatedAt,
		)

		if err != nil {
			_ = rows.Close()
			return nil, err
		}

		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}

	return items, nil
}

func (r *repo) Delete(ctx context.Context, todoID int64) error {
	q := `DELETE FROM todos WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, todoID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repo) DeleteAll(ctx context.Context) error {
	q := `DELETE FROM todos`
	res, err := r.db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
