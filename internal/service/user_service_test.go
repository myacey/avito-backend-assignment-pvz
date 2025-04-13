package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service/mocks"
)

var (
	tokenValid              = "valid"
	errMock    error        = errors.New("mock error")
	mockUser   *entity.User = &entity.User{ID: uuid.New(), Email: "mock@example.com", Password: "string", Role: entity.RoleEmployee}
)

func TestDummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	tokenSrv := mocks.NewMockTokenService(ctrl)
	srv := service.NewUserService(nil, nil, tokenSrv)
	testCases := []struct {
		name         string
		req          *request.DummyLogin
		mockBehavior func(req *request.DummyLogin)
		expResp      *response.Login
		expErr       error
	}{
		{
			name: "OK",
			req: &request.DummyLogin{
				Role: string(entity.RoleEmployee),
			},
			mockBehavior: func(req *request.DummyLogin) {
				tokenSrv.EXPECT().CreateDummyToken(req.Role).Return(tokenValid, nil)
			},
			expResp: &response.Login{tokenValid},
			expErr:  nil,
		},
		{
			name: "create dummy token err",
			req: &request.DummyLogin{
				Role: string(entity.RoleEmployee),
			},
			mockBehavior: func(req *request.DummyLogin) {
				tokenSrv.EXPECT().CreateDummyToken(req.Role).Return("", errMock)
			},
			expResp: nil,
			expErr:  apperror.NewUnauthorized(errMock.Error()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			res, err := srv.DummyLogin(context.Background(), tc.req)

			require.Equal(t, tc.expResp, res)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)

	tokenSrv := mocks.NewMockTokenService(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	srv := service.NewUserService(userRepo, nil, tokenSrv)
	testCases := []struct {
		name         string
		req          *request.Register
		mockBehavior func(req *request.Register)
		expResp      *entity.User
		expErr       error
	}{
		{
			name: "OK",
			req: &request.Register{
				Email:    mockUser.Email,
				Password: mockUser.Password,
				Role:     string(mockUser.Role),
			},
			mockBehavior: func(req *request.Register) {
				userRepo.EXPECT().CreateUser(gomock.Any(), req).Return(mockUser, nil)
			},
			expResp: mockUser,
			expErr:  nil,
		},
		{
			name: "err user already exists",
			req: &request.Register{
				Email:    mockUser.Email,
				Password: mockUser.Password,
				Role:     string(mockUser.Role),
			},
			mockBehavior: func(req *request.Register) {
				userRepo.EXPECT().CreateUser(gomock.Any(), req).Return(nil, repository.ErrUserAlreadyExists)
			},
			expResp: nil,
			expErr:  apperror.NewBadReq(repository.ErrUserAlreadyExists.Error()),
		},
		{
			name: "internal error",
			req: &request.Register{
				Email:    mockUser.Email,
				Password: mockUser.Password,
				Role:     string(mockUser.Role),
			},
			mockBehavior: func(req *request.Register) {
				userRepo.EXPECT().CreateUser(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("cant add new user", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			res, err := srv.Register(context.Background(), tc.req)

			require.Equal(t, tc.expResp, res)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	tokenSrv := mocks.NewMockTokenService(ctrl)
	userRepo := mocks.NewMockUserRepo(ctrl)

	srv := service.NewUserService(userRepo, nil, tokenSrv)
	testCases := []struct {
		name         string
		req          *request.Login
		mockBehavior func(req *request.Login)
		expResp      *response.Login
		expErr       error
	}{
		{
			name: "OK",
			req: &request.Login{
				Email:    mockUser.Email,
				Password: mockUser.Password,
			},
			mockBehavior: func(req *request.Login) {
				userRepo.EXPECT().GetUser(gomock.Any(), req).Return(mockUser, nil)
				tokenSrv.EXPECT().CreateUserToken(mockUser.ID, string(mockUser.Role)).Return(tokenValid, nil)
			},
			expResp: &response.Login{
				Token: tokenValid,
			},
			expErr: nil,
		},
		{
			name: "get user err",
			req: &request.Login{
				Email:    mockUser.Email,
				Password: mockUser.Password,
			},
			mockBehavior: func(req *request.Login) {
				userRepo.EXPECT().GetUser(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  errMock,
		},
		{
			name: "invalid password",
			req: &request.Login{
				Email:    mockUser.Email,
				Password: "invalid",
			},
			mockBehavior: func(req *request.Login) {
				userRepo.EXPECT().GetUser(gomock.Any(), req).Return(mockUser, nil)
			},
			expResp: nil,
			expErr:  apperror.NewUnauthorized("user not found"),
		},
		{
			name: "creation token err",
			req: &request.Login{
				Email:    mockUser.Email,
				Password: mockUser.Password,
			},
			mockBehavior: func(req *request.Login) {
				userRepo.EXPECT().GetUser(gomock.Any(), req).Return(mockUser, nil)
				tokenSrv.EXPECT().CreateUserToken(mockUser.ID, string(mockUser.Role)).Return("", errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to create token", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			res, err := srv.Login(context.Background(), tc.req)

			require.Equal(t, tc.expResp, res)
			require.Equal(t, tc.expErr, err)
		})
	}
}
