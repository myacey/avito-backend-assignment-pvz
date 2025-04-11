package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
)

type PvzService interface {
	CreatePvz(context.Context, *request.CreatePvz) (*entity.Pvz, error)
}

func CreatePvz(ctx *gin.Context, service PvzService) error {
	log.SetPrefix("http-server.handler.CreatePvz")
	var req request.CreatePvz
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req: " + err.Error())
	}

	resp, err := service.CreatePvz(ctx, &req)
	if err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, resp)
	return nil
}
