package repository

import db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"

type PvzRepository struct {
	queries *db.Queries
}

func NewPvzRepository(q *db.Queries) *PvzRepository {
	return &PvzRepository{q}
}
