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
	ErrReceptionInProgress  = errors.New("reception in progress")
	ErrNoOpenReceptionFound = errors.New("no in-progress reception found")
)

const (
	errReceptionInProgressConflictCode = "20001"
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
		case ok && pqErr.Code == errReceptionInProgressConflictCode:
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
