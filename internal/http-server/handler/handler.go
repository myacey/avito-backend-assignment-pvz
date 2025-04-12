package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
)

type Handler struct {
	receptionSrv ReceptionService
	pvzSrv       PvzService
	userSrv      UserService
}

func NewHandler(receptionSrv ReceptionService, pvzSrv PvzService, usrSrv UserService) *Handler {
	return &Handler{
		receptionSrv: receptionSrv,
		pvzSrv:       pvzSrv,
		userSrv:      usrSrv,
	}
}

func wrapCtxWithError(ctx *gin.Context, err error) {
	if httpError, ok := err.(apperror.HTTPError); ok {
		ctx.JSON(httpError.Code, response.Error{
			Code:      httpError.Code,
			Message:   httpError.Message,
			RequestId: ctx.GetHeader("X-Request-Id"),
		})

		if httpError.Code == http.StatusInternalServerError {
			log.Printf("internal error: %v | %v", httpError.Message, httpError.DebugError)
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, response.Error{
			Code:      http.StatusInternalServerError,
			Message:   err.Error(),
			RequestId: ctx.GetHeader("X-Request-Id"),
		})
	}
	ctx.Set("Retry-After", 10)
}
