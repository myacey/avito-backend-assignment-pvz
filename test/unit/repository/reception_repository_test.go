package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository/mocks"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
	"github.com/stretchr/testify/require"
)

var (
	pvz       = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now(), City: entity.CITY_MOSCOW}
	reception = &entity.Reception{ID: uuid.New(), DateTime: time.Now(), PvzID: pvz.ID, Status: entity.STATUS_IN_PROGRESS}
	product   = &entity.Product{ID: uuid.New(), DateTime: time.Now(), Type: entity.PRODUCT_TYPE_CLOTHES, ReceptionID: reception.ID}

	pvz1        = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now(), City: entity.CITY_MOSCOW}
	pvz2        = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Now().AddDate(0, 0, -1), City: entity.CITY_KAZAN}
	reception1  = &entity.Reception{ID: uuid.New(), DateTime: time.Now().AddDate(0, 0, -1), PvzID: pvz1.ID, Status: entity.STATUS_FINISHED}
	reception11 = &entity.Reception{ID: uuid.New(), DateTime: time.Now(), PvzID: pvz1.ID, Status: entity.STATUS_IN_PROGRESS}
	reception2  = &entity.Reception{ID: uuid.New(), DateTime: time.Now(), PvzID: pvz2.ID, Status: entity.STATUS_IN_PROGRESS}
)

func TestGetLastOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          uuid.UUID
		mockBehavior func(req uuid.UUID)
		expRes       *entity.Reception
		expErr       error
	}{
		{
			name: "ok",
			req:  pvz.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().GetOpenReceptionByPvzID(gomock.Any(), req).Return(db.Reception{
					ID:       reception.ID,
					DateTime: reception.DateTime,
					PvzID:    reception.PvzID,
					Status:   reception.Status,
				}, nil)
			},
			expRes: reception,
			expErr: nil,
		},
		{
			name: "no reception found",
			req:  pvz.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().GetOpenReceptionByPvzID(gomock.Any(), req).Return(db.Reception{}, sql.ErrNoRows)
			},
			expRes: nil,
			expErr: repository.ErrNoOpenReceptionFound,
		},
		{
			name: "unk error",
			req:  pvz.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().GetOpenReceptionByPvzID(gomock.Any(), req).Return(db.Reception{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.GetLastOpenReception(context.Background(), tc.req)

		require.Equal(t, tc.expRes, res)
		require.Equal(t, tc.expErr, err)
	}
}

func TestCreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          *request.CreateReception
		mockBehavior func(req *request.CreateReception)
		expRes       *entity.Reception
		expErr       error
	}{
		{
			name: "ok",
			req: &request.CreateReception{
				PvzID: pvz.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				queries.EXPECT().CreateReception(gomock.Any(), gomock.Any()).Return(db.Reception{
					ID:       reception.ID,
					DateTime: time.Now(),
					PvzID:    req.PvzID,
					Status:   entity.STATUS_IN_PROGRESS,
				}, nil)
			},
			expRes: reception,
			expErr: nil,
		},
		{
			name: "unk err",
			req: &request.CreateReception{
				PvzID: pvz.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				queries.EXPECT().CreateReception(gomock.Any(), gomock.Any()).Return(db.Reception{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
		{
			name: "crate reception conflict",
			req: &request.CreateReception{
				PvzID: pvz.ID,
			},
			mockBehavior: func(req *request.CreateReception) {
				queries.EXPECT().CreateReception(gomock.Any(), gomock.Any()).Return(db.Reception{}, &pq.Error{Code: "20002"})
			},
			expRes: nil,
			expErr: repository.ErrReceptionInProgress,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.CreateReception(context.Background(), tc.req)

		if res != nil {
			require.Equal(t, tc.expRes.ID, res.ID)
			require.Equal(t, tc.expRes.PvzID, res.PvzID)
			require.Equal(t, tc.expRes.Status, res.Status)
			require.WithinDuration(t, tc.expRes.DateTime, res.DateTime, time.Second)
		} else {
			require.Equal(t, tc.expRes, res)
		}

		require.Equal(t, tc.expErr, err)
	}
}

func TestAddProductToReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          *request.AddProduct
		receptionID  uuid.UUID
		mockBehavior func(req *request.AddProduct)
		expRes       *entity.Product
		expErr       error
	}{
		{
			name: "ok",
			req: &request.AddProduct{
				Type:  string(entity.PRODUCT_TYPE_CLOTHES),
				PvzID: pvz.ID,
			},
			receptionID: reception.ID,
			mockBehavior: func(req *request.AddProduct) {
				queries.EXPECT().AddProductToReception(gomock.Any(), gomock.Any()).Return(db.Product{
					ID:          product.ID,
					DateTime:    product.DateTime,
					Type:        product.Type,
					ReceptionID: product.ReceptionID,
				}, nil)
			},
			expRes: product,
			expErr: nil,
		},
		{
			name: "err other reception in progress conflict",
			req: &request.AddProduct{
				Type:  string(entity.PRODUCT_TYPE_CLOTHES),
				PvzID: pvz.ID,
			},
			receptionID: reception.ID,
			mockBehavior: func(req *request.AddProduct) {
				queries.EXPECT().AddProductToReception(gomock.Any(), gomock.Any()).Return(db.Product{}, &pq.Error{Code: "20001"})
			},
			expRes: nil,
			expErr: repository.ErrReceptionInProgress,
		},
		{
			name: "unk err",
			req: &request.AddProduct{
				Type:  string(entity.PRODUCT_TYPE_CLOTHES),
				PvzID: pvz.ID,
			},
			receptionID: reception.ID,
			mockBehavior: func(req *request.AddProduct) {
				queries.EXPECT().AddProductToReception(gomock.Any(), gomock.Any()).Return(db.Product{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.AddProductToReception(context.Background(), tc.req, tc.receptionID)

		if res != nil {
			require.Equal(t, tc.expRes.ID, res.ID)
			require.Equal(t, tc.expRes.ReceptionID, res.ReceptionID)
			require.Equal(t, tc.expRes.Type, res.Type)
			require.WithinDuration(t, tc.expRes.DateTime, res.DateTime, time.Second)
		} else {
			require.Equal(t, tc.expRes, res)
		}

		require.Equal(t, tc.expErr, err)
	}
}

func TestFinishReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          uuid.UUID
		mockBehavior func(req uuid.UUID)
		expRes       *entity.Reception
		expErr       error
	}{
		{
			name: "ok",
			req:  pvz.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().FinishReception(gomock.Any(), req).Return(db.Reception{
					ID:       reception.ID,
					DateTime: reception.DateTime,
					PvzID:    reception.PvzID,
					Status:   entity.STATUS_FINISHED,
				}, nil)
			},
			expRes: &entity.Reception{
				ID:       reception.ID,
				DateTime: reception.DateTime,
				PvzID:    reception.PvzID,
				Status:   entity.STATUS_FINISHED,
			},
			expErr: nil,
		},
		{
			name: "err no reception found",
			req:  pvz.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().FinishReception(gomock.Any(), req).Return(db.Reception{}, sql.ErrNoRows)
			},
			expRes: nil,
			expErr: repository.ErrNoOpenReceptionFound,
		},
		{
			name: "unk err",
			req:  pvz.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().FinishReception(gomock.Any(), req).Return(db.Reception{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}
	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.FinishReception(context.Background(), tc.req)

		if res != nil {
			require.Equal(t, tc.expRes.ID, res.ID)
			require.Equal(t, tc.expRes.PvzID, res.PvzID)
			require.Equal(t, tc.expRes.Status, res.Status)
			require.WithinDuration(t, tc.expRes.DateTime, res.DateTime, time.Second)
		} else {
			require.Equal(t, tc.expRes, res)
		}

		require.Equal(t, tc.expErr, err)
	}
}

func TestGetLastProductInReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          uuid.UUID
		mockBehavior func(req uuid.UUID)
		expRes       *entity.Product
		expErr       error
	}{
		{
			name: "ok",
			req:  reception.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().GetLastProductInReception(gomock.Any(), req).Return(db.Product{
					ID:          product.ID,
					DateTime:    product.DateTime,
					Type:        product.Type,
					ReceptionID: product.ReceptionID,
				}, nil)
			},
			expRes: product,
			expErr: nil,
		},
		{
			name: "no product found",
			req:  reception.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().GetLastProductInReception(gomock.Any(), req).Return(db.Product{}, sql.ErrNoRows)
			},
			expRes: nil,
			expErr: repository.ErrNoProduct,
		},
	}
	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.GetLastProductInReception(context.Background(), tc.req)

		if res != nil {
			require.Equal(t, tc.expRes.ID, res.ID)
			require.Equal(t, tc.expRes.ReceptionID, res.ReceptionID)
			require.Equal(t, tc.expRes.Type, res.Type)
			require.WithinDuration(t, tc.expRes.DateTime, res.DateTime, time.Second)
		} else {
			require.Equal(t, tc.expRes, res)
		}

		require.Equal(t, tc.expErr, err)
	}
}

func TestDeleteProductInReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          uuid.UUID
		mockBehavior func(req uuid.UUID)
		expErr       error
	}{
		{
			name: "ok",
			req:  product.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().DeleteProduct(gomock.Any(), req).Return(nil)
			},
			expErr: nil,
		},
		{
			name: "err no product",
			req:  product.ID,
			mockBehavior: func(req uuid.UUID) {
				queries.EXPECT().DeleteProduct(gomock.Any(), req).Return(sql.ErrNoRows)
			},
			expErr: repository.ErrNoProduct,
		},
	}
	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		err := repo.DeleteProductInReception(context.Background(), tc.req)

		require.Equal(t, tc.expErr, err)
	}
}

func TestSearchReceptions(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockReceptionQueries(ctrl)

	repo := repository.NewReceptionRepository(queries)
	testCases := []struct {
		name         string
		req          *request.SearchPvz
		pvzIDs       []uuid.UUID
		mockBehavior func(req *request.SearchPvz, pvzIDS []uuid.UUID)
		expRes       []*entity.Reception
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
			pvzIDs: []uuid.UUID{pvz1.ID, pvz2.ID},
			mockBehavior: func(req *request.SearchPvz, pvzIDS []uuid.UUID) {
				queries.EXPECT().SearchReceptionsByPvzsAndTime(gomock.Any(), db.SearchReceptionsByPvzsAndTimeParams{
					PvzIds:    pvzIDS,
					StartDate: req.StartDate,
					EndDate:   req.EndDate,
				}).Return([]db.Reception{
					{ID: reception1.ID, DateTime: reception1.DateTime, PvzID: reception1.PvzID, Status: reception1.Status},
					{ID: reception11.ID, DateTime: reception11.DateTime, PvzID: reception11.PvzID, Status: reception11.Status},
					{ID: reception2.ID, DateTime: reception2.DateTime, PvzID: reception2.PvzID, Status: reception2.Status},
				}, nil)
			},
			expRes: []*entity.Reception{
				{ID: reception1.ID, DateTime: reception1.DateTime, PvzID: reception1.PvzID, Status: reception1.Status},
				{ID: reception11.ID, DateTime: reception11.DateTime, PvzID: reception11.PvzID, Status: reception11.Status},
				{ID: reception2.ID, DateTime: reception2.DateTime, PvzID: reception2.PvzID, Status: reception2.Status},
			},
			expErr: nil,
		},
		{
			name: "no receptions found",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			pvzIDs: []uuid.UUID{pvz1.ID, pvz2.ID},
			mockBehavior: func(req *request.SearchPvz, pvzIDS []uuid.UUID) {
				queries.EXPECT().SearchReceptionsByPvzsAndTime(gomock.Any(), db.SearchReceptionsByPvzsAndTimeParams{
					PvzIds:    pvzIDS,
					StartDate: req.StartDate,
					EndDate:   req.EndDate,
				}).Return([]db.Reception{}, sql.ErrNoRows)
			},
			expRes: []*entity.Reception{},
			expErr: nil,
		},
		{
			name: "unk err",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     10,
			},
			pvzIDs: []uuid.UUID{pvz1.ID, pvz2.ID},
			mockBehavior: func(req *request.SearchPvz, pvzIDS []uuid.UUID) {
				queries.EXPECT().SearchReceptionsByPvzsAndTime(gomock.Any(), db.SearchReceptionsByPvzsAndTimeParams{
					PvzIds:    pvzIDS,
					StartDate: req.StartDate,
					EndDate:   req.EndDate,
				}).Return([]db.Reception{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req, tc.pvzIDs)

		res, err := repo.SearchReceptions(context.Background(), tc.req, tc.pvzIDs)

		if res != nil {
			require.Len(t, res, len(tc.expRes))

			for i := range tc.expRes {
				// Проверка всех полей, кроме времени
				require.Equal(t, tc.expRes[i].ID, res[i].ID)
				require.Equal(t, tc.expRes[i].PvzID, res[i].PvzID)
				require.Equal(t, tc.expRes[i].Status, res[i].Status)

				// Проверка времени с допуском
				require.WithinDuration(t,
					tc.expRes[i].DateTime,
					res[i].DateTime,
					time.Second, // Допустимая погрешность
					"DateTime mismatch for element %d", i,
				)
			}
		} else {
			require.Equal(t, tc.expRes, res)
		}

		require.Equal(t, tc.expErr, err)
	}
}
