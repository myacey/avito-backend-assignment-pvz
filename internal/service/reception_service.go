package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
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
	openRecepton, err := s.receptionRepo.GetLastOpenReception(ctx, req.PvzID)
	if err != nil && !errors.Is(err, repository.ErrNoOpenReceptionFound) {
		return nil, apperror.NewInternal("failed to create reception", err)
	}
	if err == nil {
		return nil, apperror.NewBadReq("can't start new reception, already in-progress: " + openRecepton.ID.String())
	}

	reception, err := s.receptionRepo.CreateReception(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReceptionInProgress):
			return nil, apperror.NewBadReq("can't start new reception, already in-progress: " + reception.ID.String())
		default:
			return nil, apperror.NewInternal("failed to create reception", err)
		}
	}

	return &response.Reception{
		ID:       reception.ID,
		DateTime: reception.DateTime,
		PvzId:    reception.PvzID,
		Status:   string(reception.Status),
	}, nil
}

func (s *ReceptionServiceImpl) AddProductToReception(ctx context.Context, req *request.AddProduct) error {
	return nil
}
