// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddProductToReception(ctx context.Context, arg AddProductToReceptionParams) (Product, error)
	CreatePVZ(ctx context.Context, arg CreatePVZParams) (Pvz, error)
	CreateReception(ctx context.Context, arg CreateReceptionParams) (Reception, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProductFromReception(ctx context.Context, id uuid.UUID) error
	FinishReception(ctx context.Context, id uuid.UUID) error
	GetOpenReceptionByPvzID(ctx context.Context, pvzID uuid.UUID) (Reception, error)
	GetProductsFromReception(ctx context.Context, receptionID uuid.NullUUID) ([]Product, error)
	GetReceptionsByPvzAndTime(ctx context.Context, arg GetReceptionsByPvzAndTimeParams) ([]Reception, error)
	GetReceptionsByTime(ctx context.Context, arg GetReceptionsByTimeParams) ([]Reception, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	SearchPvz(ctx context.Context, arg SearchPvzParams) ([]Pvz, error)
}

var _ Querier = (*Queries)(nil)
