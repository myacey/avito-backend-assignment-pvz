//go:generate mockgen -source=./handler.go -destination=./mocks/handler.go -package=mocks

package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
)

const (
	HeaderRequestID  = "X-Request-Id"
	CtxKeyRetryAfter = "Retry-After"
)

// RoleCheckerMiddleware is middleware interface
// for checking account's roles.
type RoleCheckerMiddleware interface {
	AuthMiddleware(neededRole ...entity.Role) gin.HandlerFunc
}

type Handler struct {
	receptionSrv ReceptionService
	pvzSrv       PvzService
	userSrv      UserService

	authSrv RoleCheckerMiddleware
}

func NewHandler(receptionSrv ReceptionService, pvzSrv PvzService, usrSrv UserService, autSrv RoleCheckerMiddleware) *Handler {
	return &Handler{
		receptionSrv: receptionSrv,
		pvzSrv:       pvzSrv,
		userSrv:      usrSrv,
		authSrv:      autSrv,
	}
}

func wrapCtxWithError(ctx *gin.Context, err error) {
	if httpError, ok := err.(apperror.HTTPError); ok {
		ctx.JSON(httpError.Code, response.Error{
			Code:      httpError.Code,
			Message:   httpError.Message,
			RequestId: ctx.GetHeader(HeaderRequestID),
		})

		if httpError.Code == http.StatusInternalServerError {
			log.Printf("internal error: %v | %v", httpError.Message, httpError.DebugError)
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, response.Error{
			Code:      http.StatusInternalServerError,
			Message:   err.Error(),
			RequestId: ctx.GetHeader(HeaderRequestID),
		})
	}
	ctx.Set(CtxKeyRetryAfter, 10)
}
