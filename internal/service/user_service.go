package service

import (
	"context"
	"database/sql"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository

	conn *sql.DB
}

func NewUserService(repo repository.UserRepository, conn *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
		conn: conn,
	}
}

func (s *UserServiceImpl) DummyLogin(context.Context, *request.DummyLogin) (*response.Login, error) {
	return nil, nil
}

func (s *UserServiceImpl) Register(context.Context, *request.Register) (*response.Login, error) {
	return nil, nil
}

func (s *UserServiceImpl) Login(context.Context, *request.Login) (*response.Login, error) {
	return nil, nil
}
