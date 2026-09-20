package usecase

import (
	"context"
	"errors"
	"time"
	"todo/internal/domain/converter"
	"todo/internal/domain/model"
	"todo/internal/infra/password"
	"todo/internal/repo"
	"uuid"
)

var (
	ErrCredentials = errors.New("invalid credentials")
)

type User interface {
	Register(ctx context.Context, req model.RegisterRequest) (string, error)
	Login(ctx context.Context, req model.LoginRequest) (string, error)
	Update(ctx context.Context, req model.UserUpdateRequest) error
	Delete(ctx context.Context, req model.UserDeleteRequest) error
	Get(ctx context.Context, userID string) (*model.UserResponse, error)
}

type user struct {
	repoUser repo.User
}

func NewUser(repoUser repo.User) User {
	return &user{repoUser: repoUser}
}

func (u *user) Register(ctx context.Context, req model.RegisterRequest) (string, error) {
	hash, err := password.Hash(req.Password)
	if err != nil {
		return "", err
	}

	id := uuid.NewV7().String()
	now := time.Now().UTC()
	if err = u.repoUser.Create(ctx, repo.UserCreateParams{
		ID:        id,
		Email:     req.Email,
		Name:      req.Name,
		Username:  req.Username,
		Password:  hash,
		CreatedAt: now,
		UpdatedAt: now,
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (u *user) Login(ctx context.Context, req model.LoginRequest) (string, error) {
	i, err := u.repoUser.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}

	if !password.Check(i.Password, req.Password) {
		return "", ErrCredentials
	}

	return i.ID, nil
}

func (u *user) Update(ctx context.Context, req model.UserUpdateRequest) error {
	i, err := u.repoUser.Get(ctx, req.UserID)
	if err != nil {
		return err
	}

	if req.Password != nil {
		if req.OldPassword == nil {
			return ErrCredentials
		}

		if !password.Check(i.Password, *req.OldPassword) {
			return ErrCredentials
		}

		hash, err := password.Hash(*req.Password)
		if err != nil {
			return err
		}

		req.Password = &hash
	}

	return u.repoUser.Update(ctx, repo.UserUpdateParams{
		ID:        req.UserID,
		Email:     req.Email,
		Name:      req.Name,
		Username:  req.Username,
		Password:  req.Password,
		UpdatedAt: time.Now().UTC(),
	})
}

func (u *user) Delete(ctx context.Context, req model.UserDeleteRequest) error {
	i, err := u.repoUser.Get(ctx, req.UserID)
	if err != nil {
		return err
	}

	if !password.Check(i.Password, req.Password) {
		return ErrCredentials
	}

	return u.repoUser.Delete(ctx, req.UserID)
}

func (u *user) Get(ctx context.Context, userID string) (*model.UserResponse, error) {
	i, err := u.repoUser.Get(ctx, userID)
	return converter.UserToResponse(i), err
}
