package service

import (
	"context"
	"database/sql"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

type PvzServiceImpl struct {
	repo repository.PvzRepository

	conn *sql.DB
}

func NewPvzService(repo repository.PvzRepository, conn *sql.DB) *PvzServiceImpl {
	return &PvzServiceImpl{
		repo: repo,
		conn: conn,
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
		return nil, apperror.NewBadReq("invalid req")
	}

	return resp, err
}
