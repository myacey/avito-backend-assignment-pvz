package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
)

var (
	ErrReceptionInProgress  = errors.New("other reception in progress")
	ErrNoOpenReceptionFound = errors.New("no in-progress reception found")
	ErrNoProduct            = errors.New("no product in reception")
)

const (
	errReceptionInProgressConflictCode = "20001"
	errAddReceptionToFinishedReception = "20002"
)

type ReceptionRepository struct {
	queries *db.Queries
}

func NewReceptionRepository(q *db.Queries) *ReceptionRepository {
	return &ReceptionRepository{q}
}

func (r *ReceptionRepository) GetLastOpenReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	res, err := r.queries.GetOpenReceptionByPvzID(ctx, pvzID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoOpenReceptionFound
		default:
			return nil, err
		}
	}

	return &entity.Reception{
		ID:       res.ID,
		DateTime: res.DateTime,
		PvzID:    res.PvzID,
		Status:   res.Status,
	}, nil
}

func (r *ReceptionRepository) CreateReception(ctx context.Context, req *request.CreateReception) (*entity.Reception, error) {
	arg := db.CreateReceptionParams{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PvzID:    req.PvzID,
	}

	res, err := r.queries.CreateReception(ctx, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		switch {
		case ok && pqErr.Code == errAddReceptionToFinishedReception:
			return nil, ErrReceptionInProgress
		default:
			return nil, err
		}
	}

	return &entity.Reception{
		ID:       res.ID,
		DateTime: res.DateTime,
		PvzID:    req.PvzID,
		Status:   res.Status,
	}, nil
}

func (r *ReceptionRepository) AddProductToReception(ctx context.Context, req *request.AddProduct, receptionID uuid.UUID) (*entity.Product, error) {
	arg := db.AddProductToReceptionParams{
		ID:          uuid.New(),
		Type:        entity.ProductType(req.Type),
		ReceptionID: receptionID,
	}

	res, err := r.queries.AddProductToReception(ctx, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		switch {
		case ok && pqErr.Code == errReceptionInProgressConflictCode:
			return nil, ErrReceptionInProgress
		default:
			return nil, err
		}
	}

	return &entity.Product{
		ID:          res.ID,
		DateTime:    res.DateTime,
		Type:        entity.ProductType(res.Type),
		ReceptionID: res.ReceptionID,
	}, nil
}

func (r *ReceptionRepository) SearchReceptions(ctx context.Context, req *request.SearchPvz, pvzIDs []uuid.UUID) ([]*entity.Reception, error) {
	arg := db.SearchReceptionsByPvzsAndTimeParams{
		PvzIds:    pvzIDs,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	res, err := r.queries.SearchReceptionsByPvzsAndTime(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*entity.Reception{}, nil
		default:
			return nil, err
		}
	}

	ans := make([]*entity.Reception, len(res))
	for i, r := range res {
		ans[i] = &entity.Reception{
			ID:       r.ID,
			DateTime: r.DateTime,
			PvzID:    r.PvzID,
			Status:   r.Status,
		}
	}

	return ans, nil
}

func (r *ReceptionRepository) FinishReception(ctx context.Context, pvzID uuid.UUID) (*entity.Reception, error) {
	res, err := r.queries.FinishReception(ctx, pvzID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoOpenReceptionFound
		default:
			return nil, err
		}
	}

	return &entity.Reception{
		ID:       res.ID,
		DateTime: res.DateTime,
		PvzID:    res.PvzID,
		Status:   res.Status,
	}, nil
}

func (r *ReceptionRepository) GetLastProductInReception(ctx context.Context, receptionID uuid.UUID) (*entity.Product, error) {
	res, err := r.queries.GetLastProductInReception(ctx, receptionID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoProduct
	}

	return &entity.Product{
		ID:          res.ID,
		DateTime:    res.DateTime,
		Type:        entity.ProductType(res.Type),
		ReceptionID: res.ReceptionID,
	}, nil
}

func (r *ReceptionRepository) DeleteProductInReception(ctx context.Context, productID uuid.UUID) error {
	err := r.queries.DeleteProduct(ctx, productID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoProduct
	}

	return nil
}
