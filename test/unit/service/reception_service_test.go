package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service/mocks"
	"github.com/stretchr/testify/require"
)

var (
	pvz1 *entity.Pvz = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now().AddDate(0, 0, -3)}
	pvz2 *entity.Pvz = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now().AddDate(0, 0, -1)}
	pvz3 *entity.Pvz = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now().AddDate(0, 0, 0)}

	pvzs []*entity.Pvz = []*entity.Pvz{pvz1, pvz2, pvz3}

	reception1 *entity.Reception = &entity.Reception{ID: uuid.New(), DateTime: time.Now().AddDate(0, 0, -3), PvzID: pvz1.ID, Status: entity.STATUS_FINISHED}
	reception2 *entity.Reception = &entity.Reception{ID: uuid.New(), DateTime: time.Now().AddDate(0, 0, -1), PvzID: pvz2.ID, Status: entity.STATUS_FINISHED}
	reception3 *entity.Reception = &entity.Reception{ID: uuid.New(), DateTime: time.Now().AddDate(0, 0, 0), PvzID: pvz3.ID, Status: entity.STATUS_IN_PROGRESS}

	product *entity.Product = &entity.Product{ID: uuid.New(), DateTime: time.Now(), Type: entity.PRODUCT_TYPE_CLOTHES, ReceptionID: reception3.ID}
)

func TestSearchReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	receptionRepo := mocks.NewMockReceptionRepo(ctrl)
	pvzSrv := mocks.NewMockPvzFinder(ctrl)

	srv := service.NewReceptionService(receptionRepo, nil, pvzSrv)

	testCases := []struct {
		name         string
		req          *request.SearchPvz
		mockBehavior func(req *request.SearchPvz)
		expResp      []*entity.PvzWithReception
		expErr       error
	}{
		{
			name: "ok",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			mockBehavior: func(req *request.SearchPvz) {
				pvzSrv.EXPECT().SearchPvz(gomock.Any(), req).Return(pvzs, nil)
				receptionRepo.EXPECT().SearchReceptions(gomock.Any(), req, []uuid.UUID{pvz1.ID, pvz2.ID, pvz3.ID}).Return([]*entity.Reception{reception2, reception3}, nil)
			},
			expResp: []*entity.PvzWithReception{
				{
					Pvz:        pvz1,
					Receptions: nil,
				},
				{
					Pvz:        pvz2,
					Receptions: []*entity.Reception{reception2},
				},
				{
					Pvz:        pvz3,
					Receptions: []*entity.Reception{reception3},
				},
			},
			expErr: nil,
		},
		{
			name: "search pvz err",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			mockBehavior: func(req *request.SearchPvz) {
				pvzSrv.EXPECT().SearchPvz(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  errMock,
		},
		{
			name: "search receptions err",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			mockBehavior: func(req *request.SearchPvz) {
				pvzSrv.EXPECT().SearchPvz(gomock.Any(), req).Return(pvzs, nil)
				receptionRepo.EXPECT().SearchReceptions(gomock.Any(), req, []uuid.UUID{pvz1.ID, pvz2.ID, pvz3.ID}).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  errMock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			res, err := srv.SearchReceptions(context.Background(), tc.req)

			require.Equal(t, tc.expResp, res)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestFinishReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	receptionRepo := mocks.NewMockReceptionRepo(ctrl)

	srv := service.NewReceptionService(receptionRepo, nil, nil)

	testCases := []struct {
		name         string
		req          uuid.UUID
		mockBehavior func(req uuid.UUID)
		expResp      *entity.Reception
		expErr       error
	}{
		{
			name: "ok",
			req:  pvz1.ID,
			mockBehavior: func(req uuid.UUID) {
				receptionRepo.EXPECT().FinishReception(gomock.Any(), req).Return(reception1, nil)
			},
			expResp: reception1,
			expErr:  nil,
		},
		{
			name: "no open reception err",
			req:  pvz1.ID,
			mockBehavior: func(req uuid.UUID) {
				receptionRepo.EXPECT().FinishReception(gomock.Any(), req).Return(nil, repository.ErrNoOpenReceptionFound)
			},
			expResp: nil,
			expErr:  apperror.NewBadReq(repository.ErrNoOpenReceptionFound.Error()),
		},
		{
			name: "finish reception unk err",
			req:  pvz1.ID,
			mockBehavior: func(req uuid.UUID) {
				receptionRepo.EXPECT().FinishReception(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to close last reception", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			res, err := srv.FinishReception(context.Background(), tc.req)

			require.Equal(t, tc.expResp, res)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestDeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	dbConn, txMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	receptionRepo := mocks.NewMockReceptionRepo(ctrl)

	srv := service.NewReceptionService(receptionRepo, dbConn, nil)

	testCases := []struct {
		name         string
		req          uuid.UUID
		mockBehavior func(req uuid.UUID)
		expErr       error
	}{
		{
			name: "ok",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectCommit()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(reception3, nil)
				receptionRepo.EXPECT().GetLastProductInReception(gomock.Any(), reception3.ID).Return(product, nil)
				receptionRepo.EXPECT().DeleteProductInReception(gomock.Any(), product.ID).Return(nil)
			},
			expErr: nil,
		},
		{
			name: "tx err",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin().WillReturnError(errMock)
			},
			expErr: apperror.NewInternal("failed to delete last product", errMock),
		},
		{
			name: "err no open reception found",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(nil, repository.ErrNoOpenReceptionFound)
			},
			expErr: apperror.NewBadReq("no in-progress reception found"),
		},
		{
			name: "find last open reception unk error",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(nil, errMock)
			},
			expErr: apperror.NewInternal("failed to find open reception", errMock),
		},
		{
			name: "last reception is closed err",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(reception2, nil)
			},
			expErr: apperror.NewInternal(
				"failed to find open reception",
				errors.New(
					"found closed reception while looked for closed IN SQL: "+
						reception2.ID.String()),
			),
		},
		{
			name: "no product in reception err",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(reception3, nil)
				receptionRepo.EXPECT().GetLastProductInReception(gomock.Any(), reception3.ID).Return(nil, repository.ErrNoProduct)
			},
			expErr: apperror.NewBadReq(repository.ErrNoProduct.Error()),
		},
		{
			name: "get last product in reception unk err",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(reception3, nil)
				receptionRepo.EXPECT().GetLastProductInReception(gomock.Any(), reception3.ID).Return(nil, errMock)
			},
			expErr: apperror.NewInternal("failed to find product in reception", errMock),
		},
		{
			name: "no product in repository err",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(reception3, nil)
				receptionRepo.EXPECT().GetLastProductInReception(gomock.Any(), reception3.ID).Return(product, nil)
				receptionRepo.EXPECT().DeleteProductInReception(gomock.Any(), product.ID).Return(repository.ErrNoProduct)
			},
			expErr: apperror.NewInternal("failed to delete product in reception", errors.New("found to product in reception, but found it before. id: "+product.ID.String())),
		},
		{
			name: "delete product in reception unk err",
			req:  pvz3.ID,
			mockBehavior: func(req uuid.UUID) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req).Return(reception3, nil)
				receptionRepo.EXPECT().GetLastProductInReception(gomock.Any(), reception3.ID).Return(product, nil)
				receptionRepo.EXPECT().DeleteProductInReception(gomock.Any(), product.ID).Return(errMock)
			},
			expErr: apperror.NewInternal("failed to delete product in reception", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			err := srv.DeleteLastProduct(context.Background(), tc.req)

			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestCreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	dbConn, txMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	receptionRepo := mocks.NewMockReceptionRepo(ctrl)

	srv := service.NewReceptionService(receptionRepo, dbConn, nil)

	testCases := []struct {
		name         string
		req          *request.CreateReception
		expResp      *entity.Reception
		mockBehavior func(req *request.CreateReception)
		expErr       error
	}{
		{
			name: "ok",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin()
				txMock.ExpectCommit()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(nil, repository.ErrNoOpenReceptionFound)
				receptionRepo.EXPECT().CreateReception(gomock.Any(), req).Return(reception3, nil)
			},
			expResp: reception3,
			expErr:  nil,
		},
		{
			name: "tx err",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			expResp: nil,
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin().WillReturnError(errMock)
			},
			expErr: apperror.NewInternal("failed to craete reception", errMock),
		},
		{
			name: "get last reception unk err",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to create reception", errMock),
		},
		{
			name: "found open reception err",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(reception3, nil)
			},
			expResp: nil,
			expErr:  apperror.NewBadReq("can't start new reception, already in-progress: " + reception3.ID.String()),
		},
		{
			name: "found open when shouldnt",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(reception3, repository.ErrNoOpenReceptionFound)
			},
			expResp: nil,
			expErr: apperror.NewInternal(
				"failed to find open reception",
				errors.New(
					"found open reception while looked for closed IN SQL: "+
						reception3.ID.String())),
		},
		{
			name: "create reception other in progress err",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(nil, repository.ErrNoOpenReceptionFound)
				receptionRepo.EXPECT().CreateReception(gomock.Any(), req).Return(nil, repository.ErrReceptionInProgress)
			},
			expResp: nil,
			expErr:  apperror.NewBadReq("can't start new reception, already in-progress"),
		},
		{
			name: "create reception other in progress err",
			req: &request.CreateReception{
				PvzID: pvz3.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(nil, repository.ErrNoOpenReceptionFound)
				receptionRepo.EXPECT().CreateReception(gomock.Any(), req).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to create reception", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			resp, err := srv.CreateReception(context.Background(), tc.req)

			require.Equal(t, tc.expResp, resp)
			require.Equal(t, tc.expErr, err)
		})
	}
}

func TestAddProductToReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	dbConn, txMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	receptionRepo := mocks.NewMockReceptionRepo(ctrl)

	srv := service.NewReceptionService(receptionRepo, dbConn, nil)

	testCases := []struct {
		name         string
		req          *request.AddProduct
		mockBehavior func(req *request.AddProduct)
		expResp      *entity.Product
		expErr       error
	}{
		{
			name: "ok",
			req: &request.AddProduct{
				PvzID: pvz3.ID,
				Type:  string(product.Type),
			},
			mockBehavior: func(req *request.AddProduct) {
				txMock.ExpectBegin()
				txMock.ExpectCommit()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(reception3, nil)
				receptionRepo.EXPECT().AddProductToReception(gomock.Any(), req, reception3.ID).Return(product, nil)
			},
			expResp: product,
			expErr:  nil,
		},
		{
			name: "tx err",
			req: &request.AddProduct{
				PvzID: pvz3.ID,
				Type:  string(product.Type),
			},
			mockBehavior: func(req *request.AddProduct) {
				txMock.ExpectBegin().WillReturnError(errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to add product to reception", errMock),
		},
		{
			name: "err no reception",
			req: &request.AddProduct{
				PvzID: pvz3.ID,
				Type:  string(product.Type),
			},
			mockBehavior: func(req *request.AddProduct) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(nil, repository.ErrNoOpenReceptionFound)
			},
			expResp: nil,
			expErr:  apperror.NewBadReq("no in-progress reception found"),
		},
		{
			name: "get last reception unk err",
			req: &request.AddProduct{
				PvzID: pvz3.ID,
				Type:  string(product.Type),
			},
			mockBehavior: func(req *request.AddProduct) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to add product to reception", errMock),
		},
		{
			name: "err other reception in progress",
			req: &request.AddProduct{
				PvzID: pvz3.ID,
				Type:  string(product.Type),
			},
			mockBehavior: func(req *request.AddProduct) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(reception3, nil)
				receptionRepo.EXPECT().AddProductToReception(gomock.Any(), req, reception3.ID).Return(nil, repository.ErrReceptionInProgress)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to add product to reception", errors.New("tried to add product to other open reception: id:"+reception3.ID.String())),
		},
		{
			name: "add product to reception unk err",
			req: &request.AddProduct{
				PvzID: pvz3.ID,
				Type:  string(product.Type),
			},
			mockBehavior: func(req *request.AddProduct) {
				txMock.ExpectBegin()
				txMock.ExpectRollback()

				receptionRepo.EXPECT().GetLastOpenReception(gomock.Any(), req.PvzID).Return(reception3, nil)
				receptionRepo.EXPECT().AddProductToReception(gomock.Any(), req, reception3.ID).Return(nil, errMock)
			},
			expResp: nil,
			expErr:  apperror.NewInternal("failed to add product to reception", errMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.req)

			resp, err := srv.AddProductToReception(context.Background(), tc.req)

			require.Equal(t, tc.expResp, resp)
			require.Equal(t, tc.expErr, err)
		})
	}
}
