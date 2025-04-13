package handler_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/myacey/avito-backend-assignment-pvz/pkg/openapi"
	"github.com/stretchr/testify/require"
)

var (
	reception = &entity.Reception{ID: uuid.New(), DateTime: time.Date(2022, 12, 12, 12, 12, 0, 0, time.Local), PvzID: pvz.ID, Status: entity.STATUS_IN_PROGRESS}
	product   = &entity.Product{ID: uuid.New(), DateTime: reception.DateTime, Type: entity.PRODUCT_TYPE_CLOTHES, ReceptionID: reception.ID}

	start = time.Now().AddDate(0, 0, -2)
	end   = time.Now()
	page  = 1
	limit = 10
)

func TestPostProducts(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockReceptionService(ctrl)
	authSrv := mocks.NewMockRoleCheckerMiddleware(ctrl)

	handler := handler.NewHandler(service, nil, nil, authSrv)
	testCases := []struct {
		name         string
		req          interface{}
		mockBehavior func(req interface{})
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			req: &request.AddProduct{
				Type:  string(product.Type),
				PvzID: pvz.ID,
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().AddProductToReception(gomock.Any(), req).Return(product, nil)
			},
			expBody: product.ToResponse(),
			expCode: http.StatusCreated,
		},
		{
			name: "invalid req",
			req:  "invalid",
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				// service.EXPECT().AddProductToReception(gomock.Any(), req).Return(product, nil)
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "invalid product type",
			req: &request.AddProduct{
				Type:  "invalid",
				PvzID: pvz.ID,
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				// service.EXPECT().AddProductToReception(gomock.Any(), req).Return(product, nil)
			},
			// expBody: product.ToResponse(),
			expCode: http.StatusBadRequest,
		},
		{
			name: "service err",
			req: &request.AddProduct{
				Type:  string(product.Type),
				PvzID: pvz.ID,
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().AddProductToReception(gomock.Any(), req).Return(nil, errMock)
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

			r.POST("/products", handler.PostProducts)

			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusCreated {
				var resp response.Product
				body, _ := io.ReadAll(rec.Body)
				_ = json.Unmarshal(body, &resp)

				exp := tc.expBody.(*response.Product)
				require.Equal(t, *exp, resp)
			}
		})
	}
}

func TestPostPvzPvzIdDeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockReceptionService(ctrl)
	authSrv := mocks.NewMockRoleCheckerMiddleware(ctrl)

	handler := handler.NewHandler(service, nil, nil, authSrv)
	testCases := []struct {
		name         string
		pvzID        uuid.UUID
		mockBehavior func(pvzID uuid.UUID)
		expBody      interface{}
		expCode      int
	}{
		{
			name:  "ok",
			pvzID: pvz.ID,
			mockBehavior: func(pvzID uuid.UUID) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().DeleteLastProduct(gomock.Any(), pvzID).Return(nil)
			},
			expBody: "",
			expCode: http.StatusOK,
		},
		{
			name:  "service err",
			pvzID: pvz.ID,
			mockBehavior: func(pvzID uuid.UUID) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().DeleteLastProduct(gomock.Any(), pvzID).Return(errMock)
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy", nil)

			tc.mockBehavior(tc.pvzID)
			handler.PostPvzPvzIdDeleteLastProduct(ctx, tc.pvzID)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusOK {
				require.Equal(t, tc.expBody, rec.Body.String())
			}
		})
	}
}

func TestPostPvzPvzIdCloseLastReception(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockReceptionService(ctrl)
	authSrv := mocks.NewMockRoleCheckerMiddleware(ctrl)

	handler := handler.NewHandler(service, nil, nil, authSrv)
	testCases := []struct {
		name         string
		pvzID        uuid.UUID
		mockBehavior func(pvzID uuid.UUID)
		expBody      interface{}
		expCode      int
	}{
		{
			name:  "ok",
			pvzID: pvz.ID,
			mockBehavior: func(pvzID uuid.UUID) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().FinishReception(gomock.Any(), pvzID).Return(reception, nil)
			},
			expBody: reception.ToResponse(),
			expCode: http.StatusOK,
		},
		{
			name:  "service err",
			pvzID: pvz.ID,
			mockBehavior: func(pvzID uuid.UUID) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().FinishReception(gomock.Any(), pvzID).Return(nil, errMock)
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy", nil)

			tc.mockBehavior(tc.pvzID)
			handler.PostPvzPvzIdCloseLastReception(ctx, tc.pvzID)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusOK {
				expJSON, err := json.Marshal(tc.expBody)
				require.NoError(t, err)

				require.JSONEq(t, string(expJSON), rec.Body.String())
			}
		})
	}
}

func TestGetPvz(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockReceptionService(ctrl)
	authSrv := mocks.NewMockRoleCheckerMiddleware(ctrl)

	handler := handler.NewHandler(service, nil, nil, authSrv)

	receptionResp := reception.ToResponse()
	testCases := []struct {
		name         string
		params       openapi.GetPvzParams
		mockBehavior func(req openapi.GetPvzParams)
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			params: openapi.GetPvzParams{
				StartDate: &start,
				EndDate:   &end,
				Page:      &page,
				Limit:     &limit,
			},
			mockBehavior: func(req openapi.GetPvzParams) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().SearchReceptions(gomock.Any(), &request.SearchPvz{
					StartDate: *req.StartDate,
					EndDate:   *req.EndDate,
					Page:      *req.Page,
					Limit:     *req.Limit,
				}).Return([]*entity.PvzWithReception{
					{Pvz: pvz, Receptions: []*entity.Reception{reception}},
				}, nil)
			},
			expBody: []*response.PvzWithReception{
				{Pvz: pvz.ToResponse(), Receptions: []*response.Reception{receptionResp}},
			},
			expCode: http.StatusOK,
		},
		{
			name: "service err",
			params: openapi.GetPvzParams{
				StartDate: &start,
				EndDate:   &end,
				Page:      &page,
				Limit:     &limit,
			},
			mockBehavior: func(req openapi.GetPvzParams) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().SearchReceptions(gomock.Any(), &request.SearchPvz{
					StartDate: *req.StartDate,
					EndDate:   *req.EndDate,
					Page:      *req.Page,
					Limit:     *req.Limit,
				}).Return(nil, errMock)
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy", nil)

			tc.mockBehavior(tc.params)
			handler.GetPvz(ctx, tc.params)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusOK {
				expJSON, err := json.Marshal(tc.expBody)
				require.NoError(t, err)

				require.JSONEq(t, string(expJSON), rec.Body.String())
			}
		})
	}
}

func TestPostReceptions(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockReceptionService(ctrl)
	authSrv := mocks.NewMockRoleCheckerMiddleware(ctrl)

	handler := handler.NewHandler(service, nil, nil, authSrv)
	testCases := []struct {
		name         string
		req          interface{}
		mockBehavior func(req interface{})
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			req: &request.CreateReception{
				PvzID: pvz.ID,
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().CreateReception(gomock.Any(), req).Return(reception, nil)
			},
			expBody: reception.ToResponse(),
			expCode: http.StatusCreated,
		},
		{
			name: "invalid req",
			req:  "invalid",
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				// service.EXPECT().CreateReception(gomock.Any(), req).Return(reception, nil)
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "service err",
			req: &request.CreateReception{
				PvzID: pvz.ID,
			},
			mockBehavior: func(req interface{}) {
				authSrv.EXPECT().AuthMiddleware(entity.ROLE_EMPLOYEE).Return(func(ctx *gin.Context) {})
				service.EXPECT().CreateReception(gomock.Any(), req).Return(nil, errMock)
			},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			jsonBody, err := json.Marshal(tc.req)
			require.NoError(t, err)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/dummy", bytes.NewBuffer(jsonBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			tc.mockBehavior(tc.req)
			handler.PostReceptions(ctx)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusCreated {
				expJSON, err := json.Marshal(tc.expBody)
				require.NoError(t, err)

				require.JSONEq(t, string(expJSON), rec.Body.String())
			}
		})
	}
}
