package model

import "time"

type TodoCreateRequest struct {
	UserID     string     `json:"-"`
	Task       string     `json:"task" validate:"required"`
	DueDate    *time.Time `json:"dueDate" validate:"omitempty"`
	StatusID   *int64     `json:"statusID" validate:"omitempty,oneof=1 2 3 4 5"`
	PriorityID *int64     `json:"priorityID" validate:"omitempty,oneof=1 2 3 4 5"`
}

type TodoUpdateRequest struct {
	UserID     string              `json:"-"`
	ID         int64               `json:"-"`
	Task       *string             `json:"task" validate:"omitempty"`
	DueDate    Optional[time.Time] `json:"dueDate" validate:"omitempty"`
	StatusID   *int64              `json:"statusID" validate:"omitempty,oneof=1 2 3 4 5"`
	PriorityID *int64              `json:"priorityID" validate:"omitempty,oneof=1 2 3 4 5"`
}

type TodoResponse struct {
	ID         int64      `json:"id"`
	Task       string     `json:"task"`
	DueDate    *time.Time `json:"dueDate"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	StatusID   int64      `json:"statusID"`
	PriorityID int64      `json:"priorityID"`
}

type StatusResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

type PriorityResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}
