package converter

import (
	"todo/internal/domain/entity"
	"todo/internal/domain/model"
)

func UserToResponse(e *entity.User) *model.UserResponse {
	if e == nil {
		return nil
	}

	return &model.UserResponse{
		Name:      e.Name,
		Email:     e.Email,
		Username:  e.Username,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
