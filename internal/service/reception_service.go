package service

import (
	"context"
	"database/sql"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

type ReceptionServiceImpl struct {
	receptionRepo repository.ReceptionRepository

	conn *sql.DB
}

func NewReceptionService(repo repository.ReceptionRepository, conn *sql.DB) *ReceptionServiceImpl {
	return &ReceptionServiceImpl{
		receptionRepo: repo,
		conn:          conn,
	}
}

func (s *ReceptionServiceImpl) SearchReceptions(ctx context.Context, req *request.SearchPvz) (map[string]interface{}, error) {
	return nil, nil
}

func (s *ReceptionServiceImpl) CompleteReception(ctx context.Context, pvzID string) (*response.Reception, error) {
	return nil, nil
}

func (s *ReceptionServiceImpl) DeleteLastProduct(ctx context.Context, pvzID string) error {
	return nil
}

func (s *ReceptionServiceImpl) CreateReception(ctx context.Context, req *request.CreateReception) (*response.Reception, error) {
	return nil, nil
}

func (s *ReceptionServiceImpl) AddProductToReception(ctx context.Context, req *request.AddProduct) error {
	return nil
}
