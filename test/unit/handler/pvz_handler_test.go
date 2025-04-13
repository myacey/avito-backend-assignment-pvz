package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler/mocks"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/stretchr/testify/require"
)

var (
	pvz     = &entity.Pvz{ID: uuid.New(), RegistrationDate: time.Date(2025, 12, 12, 12, 12, 12, 0, time.Local), City: entity.CITY_MOSCOW}
	errMock = errors.New("mock error")
)

func TestPostPvz(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockPvzService(ctrl)
	authSrv := mocks.NewMockRoleCheckerMiddleware(ctrl)

	handler := handler.NewHandler(nil, service, nil, authSrv)
	testCases := []struct {
		name         string
		req          interface{}
		mockBehavior func(req interface{})
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			req: &request.CreatePvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             string(pvz.City),
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_MODERATOR).Return(func(ctx *gin.Context) {})
				service.EXPECT().CreatePvz(gomock.Any(), req).Return(pvz, nil)
			},
			expBody: &response.Pvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             string(pvz.City),
			},
			expCode: http.StatusCreated,
		},
		{
			name: "invalid req",
			req:  "invalid",
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_MODERATOR).Return(func(ctx *gin.Context) {})
				// service.EXPECT().CreatePvz(gomock.Any(), req).Return(pvz, nil)
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "invalid req",
			req: &request.CreatePvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             "invalid",
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_MODERATOR).Return(func(ctx *gin.Context) {})
				// service.EXPECT().CreatePvz(gomock.Any(), req).Return(pvz, nil)
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "craete pvz err",
			req: &request.CreatePvz{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             string(pvz.City),
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_MODERATOR).Return(func(ctx *gin.Context) {})
				service.EXPECT().CreatePvz(gomock.Any(), req).Return(nil, errMock)
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			r := gin.New()

			tc.mockBehavior(tc.req)

			r.POST("/pvz", handler.PostPvz)

			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusCreated {
				var resp response.Pvz
				body, _ := io.ReadAll(rec.Body)
				_ = json.Unmarshal(body, &resp)

				exp := tc.expBody.(*response.Pvz)
				require.Equal(t, exp.ID, resp.ID)
				require.Equal(t, exp.City, resp.City)
				require.WithinDuration(t, exp.RegistrationDate, resp.RegistrationDate, time.Second)
			}
		})
	}
}
