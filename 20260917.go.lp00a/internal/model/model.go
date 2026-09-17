package model

import "time"

type TodoCreateRequest struct {
	Task string `json:"task" binding:"required"`
}

type TodoUpdateRequest struct {
	Task   *string `json:"task,omitempty"`
	IsDone *bool   `json:"isDone,omitempty"`
}

type Todo struct {
	ID        int64     `json:"id"`
	Task      string    `json:"task"`
	IsDone    bool      `json:"isDone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
