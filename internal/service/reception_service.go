//go:generate mockgen -source=./reception_service.go -destination=./mocks/reception_service.go -package=mocks

package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/metrics"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

type ReceptionRepo interface {
	AddProductToReception(ctx context.Context, req *request.AddProduct, receptionID uuid.UUID) (*entity.Product, error)
	CreateReception(ctx context.Context, req *request.CreateReception) (*entity.Reception, error)
	DeleteProductInReception(ctx context.Context, productID uuid.UUID) error
	FinishReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error)
	GetLastOpenReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error)
	SearchReceptions(ctx context.Context, req *request.SearchPvz, pvzIDs []uuid.UUID) ([]*entity.Reception, error)
	GetLastProductInReception(ctx context.Context, receptionID uuid.UUID) (*entity.Product, error)
}

type PvzFinder interface {
	SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error)
}

type ReceptionServiceImpl struct {
	receptionRepo ReceptionRepo
	pvzSrv        PvzFinder

	conn *sql.DB
}

func NewReceptionService(repo ReceptionRepo, conn *sql.DB, pvzSrv PvzFinder) *ReceptionServiceImpl {
	return &ReceptionServiceImpl{
		receptionRepo: repo,
		conn:          conn,
		pvzSrv:        pvzSrv,
	}
}

func (s *ReceptionServiceImpl) SearchReceptions(ctx context.Context, req *request.SearchPvz) ([]*entity.PvzWithReception, error) {
	pvzs, err := s.pvzSrv.SearchPvz(ctx, req)
	if err != nil {
		return nil, err
	}

	pvzIDs := make([]uuid.UUID, len(pvzs))
	for i, pvz := range pvzs {
		pvzIDs[i] = pvz.ID
	}

	receptions, err := s.receptionRepo.SearchReceptions(ctx, req, pvzIDs)
	if err != nil {
		return nil, err
	}

	receptionBYPvz := make(map[uuid.UUID][]*entity.Reception)
	for _, r := range receptions {
		receptionBYPvz[r.PvzID] = append(receptionBYPvz[r.PvzID], r)
	}

	res := make([]*entity.PvzWithReception, 0, len(pvzs))
	for _, pvz := range pvzs {

		pw := &entity.PvzWithReception{
			Pvz:        pvz,
			Receptions: receptionBYPvz[pvz.ID],
		}
		res = append(res, pw)
	}

	return res, nil
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
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return apperror.NewInternal("failed to delete last product", err)
	}
	defer tx.Rollback()

	openReception, err := s.receptionRepo.GetLastOpenReception(ctx, pvzID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoOpenReceptionFound):
			return apperror.NewBadReq("no in-progress reception found")
		default:
			return apperror.NewInternal("failed to find open reception", err)
		}
	}

	if openReception.Status == entity.STATUS_FINISHED { // smt went really wrong
		return apperror.NewInternal(
			"failed to find open reception",
			errors.New(
				"found closed reception while looked for closed IN SQL: "+
					openReception.ID.String()),
		)
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
			return apperror.NewInternal("failed to delete product in reception", errors.New("found to product in reception, but found it before. id: "+lastProduct.ID.String()))
		default:
			return apperror.NewInternal("failed to delete product in reception", err)
		}
	}

	tx.Commit()
	return nil
}

func (s *ReceptionServiceImpl) CreateReception(ctx context.Context, req *request.CreateReception) (*entity.Reception, error) {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, apperror.NewInternal("failed to create reception", err)
	}
	defer tx.Rollback()

	openReception, err := s.receptionRepo.GetLastOpenReception(ctx, req.PvzID)
	if err != nil && !errors.Is(err, repository.ErrNoOpenReceptionFound) {
		return nil, apperror.NewInternal("failed to create reception", err)
	}
	if err == nil {
		return nil, apperror.NewBadReq("can't start new reception, already in-progress: " + openReception.ID.String())
	}

	if openReception != nil && openReception.Status == entity.STATUS_IN_PROGRESS { // smt went really wrong
		return nil, apperror.NewInternal(
			"failed to find open reception",
			errors.New(
				"found open reception while looked for closed IN SQL: "+
					openReception.ID.String()),
		)
	}

	reception, err := s.receptionRepo.CreateReception(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReceptionInProgress):
			return nil, apperror.NewBadReq("can't start new reception, already in-progress")
		default:
			return nil, apperror.NewInternal("failed to create reception", err)
		}
	}

	tx.Commit()
	metrics.CreateReception()
	return reception, nil
}

func (s *ReceptionServiceImpl) AddProductToReception(ctx context.Context, req *request.AddProduct) (*entity.Product, error) {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, apperror.NewInternal("failed to add product to reception", err)
	}
	defer tx.Rollback()

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
			return nil, apperror.NewInternal("failed to add product to reception", errors.New("tried to add product to other open reception: id:"+openReception.ID.String()))
		default:
			return nil, apperror.NewInternal("failed to add product to reception", err)
		}
	}

	tx.Commit()
	metrics.AddProduct()
	return res, nil
}
