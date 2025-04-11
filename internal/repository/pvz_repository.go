package repository

import (
	"context"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
)

type PvzRepository struct {
	queries *db.Queries
}

func NewPvzRepository(q *db.Queries) *PvzRepository {
	return &PvzRepository{q}
}

func (r *PvzRepository) CreatePvz(ctx context.Context, req *request.CreatePvz) (*entity.Pvz, error) {
	arg := db.CreatePVZParams{
		ID:               req.ID,
		RegistrationDate: req.RegistrationDate,
		City:             entity.City(req.City),
	}

	res, err := r.queries.CreatePVZ(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &entity.Pvz{
		ID:               req.ID,
		RegistrationDate: res.RegistrationDate,
		City:             string(res.City),
	}, nil
}
