package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler/mocks"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
)

var mockuser = &entity.User{ID: uuid.New(), Email: "mock@example.com", Password: "mockpassword", Role: entity.RoleEmployee}

func TestPostDummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockUserService(ctrl)

	handler := handler.NewHandler(nil, nil, service, nil)
	testCases := []struct {
		name         string
		req          interface{}
		mockBehavior func(req interface{})
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			req: &request.DummyLogin{
				Role: string(entity.RoleEmployee),
			},
			mockBehavior: func(req interface{}) {
				service.EXPECT().DummyLogin(gomock.Any(), req).Return(&response.Login{Token: "valid"}, nil)
			},
			expBody: &response.Login{Token: "valid"},
			expCode: http.StatusOK,
		},
		{
			name: "bad req",
			req:  "nivalid",
			mockBehavior: func(req interface{}) {
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "invalid role",
			req: &request.DummyLogin{
				Role: "invalid",
			},
			mockBehavior: func(req interface{}) {
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "service err",
			req: &request.DummyLogin{
				Role: string(entity.RoleEmployee),
			},
			mockBehavior: func(req interface{}) {
				service.EXPECT().DummyLogin(gomock.Any(), req).Return(nil, errMock)
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

			r.POST("/dummyLogin", handler.PostDummyLogin)

			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusOK {
				var resp response.Login
				body, _ := io.ReadAll(rec.Body)
				_ = json.Unmarshal(body, &resp)

				exp := tc.expBody.(*response.Login)
				require.Equal(t, exp.Token, resp.Token)
			}
		})
	}
}

func TestPostRegister(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockUserService(ctrl)

	handler := handler.NewHandler(nil, nil, service, nil)
	testCases := []struct {
		name         string
		req          interface{}
		mockBehavior func(req interface{})
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			req: &request.Register{
				Email:    mockuser.Email,
				Password: mockuser.Password,
				Role:     string(mockuser.Role),
			},
			mockBehavior: func(req interface{}) {
				service.EXPECT().Register(gomock.Any(), req).Return(mockuser, nil)
			},
			expBody: &response.User{mockuser.ID, mockuser.Email, string(mockuser.Role)},
			expCode: http.StatusCreated,
		},
		{
			name: "invalid req",
			req:  "invalid",
			mockBehavior: func(req interface{}) {
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "invalid role",
			req: &request.Register{
				Email:    mockuser.Email,
				Password: mockuser.Password,
				Role:     "invalid",
			},
			mockBehavior: func(req interface{}) {
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "service err",
			req: &request.Register{
				Email:    mockuser.Email,
				Password: mockuser.Password,
				Role:     string(entity.RoleEmployee),
			},
			mockBehavior: func(req interface{}) {
				service.EXPECT().Register(gomock.Any(), req).Return(nil, errMock)
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

			r.POST("/register", handler.PostRegister)

			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusOK {
				var resp response.User
				body, _ := io.ReadAll(rec.Body)
				_ = json.Unmarshal(body, &resp)

				exp := tc.expBody.(*response.User)
				require.Equal(t, exp, resp)
			}
		})
	}
}

func TestPostLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockUserService(ctrl)

	handler := handler.NewHandler(nil, nil, service, nil)
	testCases := []struct {
		name         string
		req          interface{}
		mockBehavior func(req interface{})
		expBody      interface{}
		expCode      int
	}{
		{
			name: "ok",
			req: &request.Login{
				Email:    mockuser.Email,
				Password: mockuser.Password,
			},
			mockBehavior: func(req interface{}) {
				service.EXPECT().Login(gomock.Any(), req).Return(&response.Login{Token: "valid"}, nil)
			},
			expBody: &response.Login{"valid"},
			expCode: http.StatusOK,
		},
		{
			name: "bad req",
			req:  "invalid",
			mockBehavior: func(req interface{}) {
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "service err",
			req: &request.Login{
				Email:    mockuser.Email,
				Password: mockuser.Password,
			},
			mockBehavior: func(req interface{}) {
				service.EXPECT().Login(gomock.Any(), req).Return(nil, errMock)
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

			r.POST("/login", handler.PostLogin)

			body, _ := json.Marshal(tc.req)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(rec, req)

			require.Equal(t, tc.expCode, rec.Code)

			if tc.expCode == http.StatusOK {
				var resp response.Login
				body, _ := io.ReadAll(rec.Body)
				_ = json.Unmarshal(body, &resp)

				exp := tc.expBody.(*response.Login)
				require.Equal(t, *exp, resp)
			}
		})
	}
}
