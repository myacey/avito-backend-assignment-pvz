package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
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

func (s *ReceptionServiceImpl) FinishReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	res, err := s.receptionRepo.FinishReception(ctx, pvzID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoOpenReceptionFound):
			return nil, apperror.NewBadReq(err.Error())
		default:
			return nil, apperror.NewInternal("failed to close last reception", err)
		}
	}

	return res, nil
}

func (s *ReceptionServiceImpl) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	openReception, err := s.receptionRepo.GetLastOpenReception(ctx, pvzID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoOpenReceptionFound):
			return apperror.NewBadReq("no in-progress reception found")
		default:
			return apperror.NewInternal("failed to find open reception", err)
		}
	}

	lastProduct, err := s.receptionRepo.GetLastProductInReception(ctx, openReception.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoProduct):
			return apperror.NewBadReq(err.Error())
		default:
			return apperror.NewInternal("failed to find product in reception", err)
		}
	}

	err = s.receptionRepo.DeleteProductInReception(ctx, lastProduct.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoProduct):
			return apperror.NewBadReq(err.Error())
		default:
			return apperror.NewInternal("failed to delete product in reception", err)
		}
	}
	return nil
}

func (s *ReceptionServiceImpl) CreateReception(ctx context.Context, req *request.CreateReception) (*entity.Reception, error) {
	openReception, err := s.receptionRepo.GetLastOpenReception(ctx, req.PvzID)
	if err != nil && !errors.Is(err, repository.ErrNoOpenReceptionFound) {
		return nil, apperror.NewInternal("failed to create reception", err)
	}
	if err == nil {
		return nil, apperror.NewBadReq("can't start new reception, already in-progress: " + openReception.ID.String())
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

	return reception, nil
}

func (s *ReceptionServiceImpl) AddProductToReception(ctx context.Context, req *request.AddProduct) (*entity.Product, error) {
	openReception, err := s.receptionRepo.GetLastOpenReception(ctx, req.PvzID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoOpenReceptionFound):
			return nil, apperror.NewBadReq("no in-progress reception found")
		default:
			return nil, apperror.NewInternal("failed to add product to reception", err)
		}
	}

	res, err := s.receptionRepo.AddProductToReception(ctx, req, openReception.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReceptionInProgress):
			return nil, apperror.NewBadReq("can't start new reception, other already in progress")
		default:
			return nil, apperror.NewInternal("failed to add product to reception", err)
		}
	}

	return res, nil
}
