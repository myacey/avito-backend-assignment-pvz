package repository

import db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"

type ReceptionRepository struct {
	queries *db.Queries
}

func NewReceptionRepository(q *db.Queries) *ReceptionRepository {
	return &ReceptionRepository{q}
}
