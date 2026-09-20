package converter

import (
	"todo/internal/domain/entity"
	"todo/internal/domain/model"
)

func TodoToResponse(e *entity.Todo) *model.TodoResponse {
	if e == nil {
		return nil
	}

	return &model.TodoResponse{
		ID:         e.ID,
		Task:       e.Task,
		DueDate:    e.DueDate,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		StatusID:   e.StatusID,
		PriorityID: e.PriorityID,
	}
}

func StatusToResponse(e *entity.Status) *model.StatusResponse {
	if e == nil {
		return nil
	}

	return &model.StatusResponse{
		ID:    e.ID,
		Name:  e.Name,
		Emoji: e.Emoji,
	}
}

func PriorityToResponse(e *entity.Priority) *model.PriorityResponse {
	if e == nil {
		return nil
	}

	return &model.PriorityResponse{
		ID:    e.ID,
		Name:  e.Name,
		Emoji: e.Emoji,
	}
}
