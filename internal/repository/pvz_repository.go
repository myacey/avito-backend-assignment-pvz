//go:generate mockgen -source=./pvz_repository.go -destination=mocks/pvz_repository.go -package=mocks

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
)

var ErrPvzAlreadyExists = errors.New("pvz already exists")

type PvzQueries interface {
	SearchPVZ(ctx context.Context, arg db.SearchPVZParams) ([]db.Pvz, error)
	CreatePVZ(ctx context.Context, arg db.CreatePVZParams) (db.Pvz, error)
}

type PvzRepository struct {
	queries PvzQueries
}

func NewPvzRepository(q PvzQueries) *PvzRepository {
	return &PvzRepository{q}
}

func (r *PvzRepository) SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error) {
	arg := db.SearchPVZParams{
		Offset: (int32(req.Page) - 1) * int32(req.Limit),
		Limit:  int32(req.Limit),
	}

	res, err := r.queries.SearchPVZ(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*entity.Pvz{}, nil
		default:
			return nil, err
		}
	}

	pvz := make([]*entity.Pvz, len(res))
	for i, r := range res {
		pvz[i] = &entity.Pvz{
			ID:               r.ID,
			RegistrationDate: r.RegistrationDate,
			City:             r.City,
		}
	}

	return pvz, nil
}

func (r *PvzRepository) CreatePvz(ctx context.Context, req *request.CreatePvz) (*entity.Pvz, error) {
	arg := db.CreatePVZParams{
		ID:               req.ID,
		RegistrationDate: req.RegistrationDate,
		City:             entity.City(req.City),
	}

	res, err := r.queries.CreatePVZ(ctx, arg)
	if err != nil {
		switch {
		case isUniqueViolation(err):
			return nil, ErrPvzAlreadyExists
		default:
			return nil, err
		}
	}

	return &entity.Pvz{
		ID:               req.ID,
		RegistrationDate: res.RegistrationDate,
		City:             res.City,
	}, nil
}
