package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/vanh01/api-rate-limiting/internal/model"
)

type (
	UserUsecase interface {
		GetById(ctx context.Context, id uuid.UUID) (*model.User, error)
	}
)

type (
	UserRepo interface {
		GetById(ctx context.Context, id uuid.UUID) (*model.User, error)
	}
)
