package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

type TokenService interface {
	CraeteDummyToken(role string) (string, error)
	CraeteUserToken(id uuid.UUID, role string) (string, error)
	VerifyToken(tokenStr string) (map[string]interface{}, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository

	conn *sql.DB

	tokenSrv TokenService
}

func NewUserService(repo repository.UserRepository, conn *sql.DB, tokenSrv TokenService) *UserServiceImpl {
	return &UserServiceImpl{
		repo:     repo,
		conn:     conn,
		tokenSrv: tokenSrv,
	}
}

func (s *UserServiceImpl) DummyLogin(ctx context.Context, req *request.DummyLogin) (*response.Login, error) {
	tokenStr, err := s.tokenSrv.CraeteDummyToken(req.Role)
	if err != nil {
		return nil, apperror.NewUnauthorized(err.Error())
	}

	return &response.Login{
		Token: tokenStr,
	}, nil
}

func (s *UserServiceImpl) Register(ctx context.Context, req *request.Register) (*entity.User, error) {
	res, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			return nil, apperror.NewBadReq(err.Error())
		default:
			return nil, apperror.NewInternal("cant add new user", err)
		}
	}

	return res, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *request.Login) (*response.Login, error) {
	res, err := s.repo.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if req.Password != res.Password {
		return nil, apperror.NewUnauthorized("user not found")
	}

	tokenStr, err := s.tokenSrv.CraeteUserToken(res.ID, string(res.Role))
	if err != nil {
		return nil, apperror.NewInternal("unable to craete", err)
	}

	return &response.Login{
		Token: tokenStr,
	}, nil
}
