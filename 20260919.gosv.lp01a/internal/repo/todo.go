package repo

import (
	"context"
	"database/sql"
	"time"
	"todo/internal/domain/entity"
)

type Todo interface {
	Create(ctx context.Context, arg TodoCreateParams) (*entity.Todo, error)
	Update(ctx context.Context, arg TodoUpdateParams) (*entity.Todo, error)
	GetAll(ctx context.Context, userID string) ([]entity.Todo, error)
	Delete(ctx context.Context, userID string, todoID int64) error
	DeleteAll(ctx context.Context, userID string) error
	GetStatuses(ctx context.Context) ([]entity.Status, error)
	GetPriorities(ctx context.Context) ([]entity.Priority, error)
}

type todo struct {
	*sql.DB
}

func NewTodo(db *sql.DB) Todo {
	return &todo{db}
}

type TodoCreateParams struct {
	Task       string
	DueDate    *time.Time
	StatusID   *int64
	PriorityID *int64
	UserID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (t *todo) Create(ctx context.Context, arg TodoCreateParams) (*entity.Todo, error) {
	q := `
	INSERT INTO todos
		( task, due_date, status_id, priority_id, created_at, updated_at, user_id )
	VALUES
		( ?1, ?2, COALESCE(?3, 1), COALESCE(?4, 3), ?5, ?6, ?7 )
	RETURNING id, task, due_date, status_id, priority_id, created_at, updated_at, user_id`

	i := new(entity.Todo)
	if err := t.QueryRowContext(ctx, q,
		arg.Task,
		arg.DueDate,
		arg.StatusID,
		arg.PriorityID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
	).Scan(
		&i.ID,
		&i.Task,
		&i.DueDate,
		&i.StatusID,
		&i.PriorityID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	); err != nil {
		return nil, err
	}

	return i, nil
}

type TodoUpdateParams struct {
	ID         int64
	UserID     string
	Task       *string
	StatusID   *int64
	PriorityID *int64
	DueDate    *time.Time
	SetDueDate bool
	UpdatedAt  time.Time
}

func (u *todo) Update(ctx context.Context, arg TodoUpdateParams) (*entity.Todo, error) {
	q := `
	UPDATE todos
	SET
		task 				= COALESCE(?1, task),
		status_id 	= COALESCE(?2, status_id),
		priority_id = COALESCE(?3, priority_id),
		due_date 		= IIF(?4, ?5, due_date),
		updated_at 	= ?6
	WHERE user_id ?7 AND id = ?8
	RETURNING id, task, due_date, status_id, priority_id, created_at, updated_at, user_id`
	i := new(entity.Todo)
	if err := u.
		QueryRowContext(ctx, q,
			arg.Task,
			arg.StatusID,
			arg.PriorityID,
			arg.SetDueDate, arg.DueDate,
			arg.UpdatedAt,
			arg.UserID,
			arg.ID,
		).
		Scan(
			&i.ID,
			&i.Task,
			&i.DueDate,
			&i.StatusID,
			&i.PriorityID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
		return nil, err
	}
	return i, nil
}

func (t *todo) GetAll(ctx context.Context, userID string) ([]entity.Todo, error) {
	q := `
	SELECT id, task, due_date, status_id, priority_id, created_at, updated_at, user_id
	FROM todos
	WHERE user_id = ?
	ORDER BY priority_id ASC, status_id ASC, due_date ASC, updated_at DESC`
	rows, err := t.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.Todo
	for rows.Next() {
		var i entity.Todo
		if err := rows.Scan(
			&i.ID,
			&i.Task,
			&i.DueDate,
			&i.StatusID,
			&i.PriorityID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
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

func (t *todo) Delete(ctx context.Context, userID string, todoID int64) error {
	q := `DELETE FROM todos WHERE user_id = ? AND id = ?`
	result, err := t.ExecContext(ctx, q, userID, todoID)
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

func (t *todo) DeleteAll(ctx context.Context, userID string) error {
	q := `DELETE FROM todos WHERE user_id = ?`
	result, err := t.ExecContext(ctx, q, userID)
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

func (t *todo) GetStatuses(ctx context.Context) ([]entity.Status, error) {
	q := "SELECT id, name, emoji FROM statuses ORDER BY id ASC"

	rows, err := t.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []entity.Status
	for rows.Next() {
		var i entity.Status
		if err = rows.Scan(
			&i.ID,
			&i.Name,
			&i.Emoji,
		); err != nil {
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

func (t *todo) GetPriorities(ctx context.Context) ([]entity.Priority, error) {
	q := "SELECT id, name, emoji FROM priorities ORDER BY id ASC"

	rows, err := t.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []entity.Priority
	for rows.Next() {
		var i entity.Priority
		if err = rows.Scan(
			&i.ID,
			&i.Name,
			&i.Emoji,
		); err != nil {
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
