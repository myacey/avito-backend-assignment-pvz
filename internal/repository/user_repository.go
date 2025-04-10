package repository

import db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q}
}
