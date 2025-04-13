package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository/mocks"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
)

func TestSearchPvz(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockPvzQueries(ctrl)

	repo := repository.NewPvzRepository(queries)
	testCases := []struct {
		name         string
		req          *request.SearchPvz
		mockBehavior func(req *request.SearchPvz)
		expRes       []*entity.Pvz
		expErr       error
	}{
		{
			name: "ok",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     1,
			},
			mockBehavior: func(req *request.SearchPvz) {
				queries.EXPECT().SearchPVZ(gomock.Any(), db.SearchPVZParams{
					Offset: (int32(req.Page) - 1) * int32(req.Limit),
					Limit:  int32(req.Limit),
				}).Return([]db.Pvz{
					{pvz1.ID, pvz1.RegistrationDate, pvz1.City},
					{pvz2.ID, pvz2.RegistrationDate, pvz2.City},
				}, nil)
			},
			expRes: []*entity.Pvz{
				{pvz1.ID, pvz1.RegistrationDate, pvz1.City},
				{pvz2.ID, pvz2.RegistrationDate, pvz2.City},
			},
			expErr: nil,
		},
		{
			name: "no pvz found",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     1,
			},
			mockBehavior: func(req *request.SearchPvz) {
				queries.EXPECT().SearchPVZ(gomock.Any(), db.SearchPVZParams{
					Offset: (int32(req.Page) - 1) * int32(req.Limit),
					Limit:  int32(req.Limit),
				}).Return([]db.Pvz{}, sql.ErrNoRows)
			},
			expRes: []*entity.Pvz{},
			expErr: nil,
		},
		{
			name: "unk err",
			req: &request.SearchPvz{
				StartDate: time.Now().AddDate(0, 0, -2),
				EndDate:   time.Now(),
				Page:      1,
				Limit:     1,
			},
			mockBehavior: func(req *request.SearchPvz) {
				queries.EXPECT().SearchPVZ(gomock.Any(), db.SearchPVZParams{
					Offset: (int32(req.Page) - 1) * int32(req.Limit),
					Limit:  int32(req.Limit),
				}).Return([]db.Pvz{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.SearchPvz(context.Background(), tc.req)

		if res != nil {
			require.Len(t, res, len(tc.expRes))

			for i := range tc.expRes {
				require.Equal(t, tc.expRes[i].ID, res[i].ID)
				require.Equal(t, tc.expRes[i].City, res[i].City)

				require.WithinDuration(t,
					tc.expRes[i].RegistrationDate,
					res[i].RegistrationDate,
					time.Second,
					"DateTime mismatch for element %d", i,
				)
			}
		} else {
			require.Equal(t, tc.expRes, res)
		}
		require.Equal(t, tc.expRes, res)
		require.Equal(t, tc.expErr, err)
	}
}

func TestCreatePvz(t *testing.T) {
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockPvzQueries(ctrl)

	repo := repository.NewPvzRepository(queries)
	testCases := []struct {
		name         string
		req          *request.CreatePvz
		mockBehavior func(req *request.CreatePvz)
		expRes       *entity.Pvz
		expErr       error
	}{
		{
			name: "ok",
			req: &request.CreatePvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             string(pvz.City),
			},
			mockBehavior: func(req *request.CreatePvz) {
				queries.EXPECT().CreatePVZ(gomock.Any(), db.CreatePVZParams{
					ID:               req.ID,
					RegistrationDate: req.RegistrationDate,
					City:             entity.City(req.City),
				}).Return(db.Pvz{
					ID:               pvz.ID,
					RegistrationDate: pvz.RegistrationDate,
					City:             pvz.City,
				}, nil)
			},
			expRes: pvz,
			expErr: nil,
		},
		{
			name: "pvzID duplicate",
			req: &request.CreatePvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             string(pvz.City),
			},
			mockBehavior: func(req *request.CreatePvz) {
				queries.EXPECT().CreatePVZ(gomock.Any(), db.CreatePVZParams{
					ID:               req.ID,
					RegistrationDate: req.RegistrationDate,
					City:             entity.City(req.City),
				}).Return(db.Pvz{}, &pq.Error{Code: "23505"})
			},
			expRes: nil,
			expErr: repository.ErrPvzAlreadyExists,
		},
		{
			name: "unk err",
			req: &request.CreatePvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             string(pvz.City),
			},
			mockBehavior: func(req *request.CreatePvz) {
				queries.EXPECT().CreatePVZ(gomock.Any(), db.CreatePVZParams{
					ID:               req.ID,
					RegistrationDate: req.RegistrationDate,
					City:             entity.City(req.City),
				}).Return(db.Pvz{}, errMock)
			},
			expRes: nil,
			expErr: errMock,
		},
	}

	for _, tc := range testCases {
		tc.mockBehavior(tc.req)

		res, err := repo.CreatePvz(context.Background(), tc.req)

		require.Equal(t, tc.expRes, res)

		require.Equal(t, tc.expRes, res)
		require.Equal(t, tc.expErr, err)
	}
}
