package usecase

import (
	"context"
	"time"
	"todo/internal/domain/converter"
	"todo/internal/domain/model"
	"todo/internal/repo"
)

type Todo interface {
	Create(ctx context.Context, req model.TodoCreateRequest) (*model.TodoResponse, error)
	Update(ctx context.Context, req model.TodoUpdateRequest) (*model.TodoResponse, error)
	GetAll(ctx context.Context, userID string) ([]model.TodoResponse, error)
	Delete(ctx context.Context, userID string, todoID int64) error
	DeleteAll(ctx context.Context, userID string) error
	GetStatuses(ctx context.Context) ([]model.StatusResponse, error)
	GetPriorities(ctx context.Context) ([]model.PriorityResponse, error)
}

type todo struct {
	repoTodo repo.Todo
}

func NewTodo(repoTodo repo.Todo) Todo {
	return &todo{repoTodo: repoTodo}
}

func (t *todo) Create(ctx context.Context, req model.TodoCreateRequest) (*model.TodoResponse, error) {
	now := time.Now().UTC()
	i, err := t.repoTodo.Create(ctx, repo.TodoCreateParams{
		Task:       req.Task,
		DueDate:    req.DueDate,
		StatusID:   req.StatusID,
		PriorityID: req.PriorityID,
		UserID:     req.UserID,
		CreatedAt:  now,
		UpdatedAt:  now,
	})

	return converter.TodoToResponse(i), err
}

func (t *todo) Update(ctx context.Context, req model.TodoUpdateRequest) (*model.TodoResponse, error) {
	i, err := t.repoTodo.Update(ctx, repo.TodoUpdateParams{
		ID:         req.ID,
		UserID:     req.UserID,
		Task:       req.Task,
		StatusID:   req.StatusID,
		PriorityID: req.PriorityID,
		DueDate:    &req.DueDate.Value,
		SetDueDate: req.DueDate.Set,
	})
	return converter.TodoToResponse(i), err
}

func (t *todo) GetAll(ctx context.Context, userID string) ([]model.TodoResponse, error) {
	items, err := t.repoTodo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	todos := make([]model.TodoResponse, 0, len(items))
	for i := range items {
		todos = append(todos, *converter.TodoToResponse(&items[i]))
	}

	return todos, nil
}

func (t *todo) Delete(ctx context.Context, userID string, todoID int64) error {
	return t.repoTodo.Delete(ctx, userID, todoID)
}

func (t *todo) DeleteAll(ctx context.Context, userID string) error {
	return t.repoTodo.DeleteAll(ctx, userID)
}

func (t *todo) GetStatuses(ctx context.Context) ([]model.StatusResponse, error) {
	items, err := t.repoTodo.GetStatuses(ctx)
	if err != nil {
		return nil, err
	}

	statuses := make([]model.StatusResponse, 0, len(items))
	for i := range items {
		statuses = append(statuses, *converter.StatusToResponse(&items[i]))
	}

	return statuses, nil
}

func (t *todo) GetPriorities(ctx context.Context) ([]model.PriorityResponse, error) {
	items, err := t.repoTodo.GetPriorities(ctx)
	if err != nil {
		return nil, err
	}

	priorities := make([]model.PriorityResponse, 0, len(items))
	for i := range items {
		priorities = append(priorities, *converter.PriorityToResponse(&items[i]))
	}

	return priorities, nil
}
