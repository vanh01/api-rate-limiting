package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vanh01/api-rate-limiting/internal/model"
)

type userRepo struct {
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (a *userRepo) GetById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return &model.User{
		ID:   id,
		Name: fmt.Sprintf("name for %s", id.String()),
	}, nil
}
