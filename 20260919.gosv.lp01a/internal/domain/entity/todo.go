package entity

import "time"

type Todo struct {
	ID         int64
	Task       string
	DueDate    *time.Time
	StatusID   int64
	PriorityID int64
	UserID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Status struct {
	ID    int64
	Name  string
	Emoji string
}

type Priority struct {
	ID    int64
	Name  string
	Emoji string
}
