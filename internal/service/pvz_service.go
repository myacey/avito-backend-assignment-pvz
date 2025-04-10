package service

import (
	"context"
	"database/sql"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
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

func (s *PvzServiceImpl) CreatePvz(ctx context.Context, req *request.CreatePvz) (*response.CreatePvz, error) {
	return nil, nil
}
