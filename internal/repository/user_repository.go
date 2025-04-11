package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q}
}

func (r *UserRepository) CreateUser(ctx context.Context, req *request.Register) (*entity.User, error) {
	arg := db.CreateUserParams{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
		Role:     entity.Role(req.Role),
	}

	res, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		switch {
		case isUniqueViolation(err):
			return nil, ErrUserAlreadyExists
		default:
			return nil, err
		}
	}

	return &entity.User{
		ID:    res.ID,
		Email: res.Email,
		Role:  res.Role,
	}, nil
}
