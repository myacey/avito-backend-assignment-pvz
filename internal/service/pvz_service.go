//go:generate mockgen -source=./pvz_service.go -destination=./mocks/pvz_service.go -package=mocks

package service

import (
	"context"
	"errors"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/metrics"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

type PvzRepo interface {
	CreatePvz(ctx context.Context, req *request.CreatePvz) (*entity.Pvz, error)
	SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error)
}

type PvzServiceImpl struct {
	repo PvzRepo
}

func NewPvzService(repo PvzRepo) *PvzServiceImpl {
	return &PvzServiceImpl{
		repo: repo,
	}
}

func (s *PvzServiceImpl) SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error) {
	res, err := s.repo.SearchPvz(ctx, req)
	if err != nil {
		return nil, apperror.NewInternal("failed to find pvz", err)
	}

	return res, nil
}

func (s *PvzServiceImpl) CreatePvz(ctx context.Context, req *request.CreatePvz) (*entity.Pvz, error) {
	resp, err := s.repo.CreatePvz(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPvzAlreadyExists):
			return nil, apperror.NewBadReq(err.Error())
		default:
			return nil, apperror.NewInternal("failed to craete repository", err)
		}
	}

	metrics.CreatePVZ()
	return resp, err
}
