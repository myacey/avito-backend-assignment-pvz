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

// PostPvz creates a new pvz with moderator auth.
func (h Handler) PostPvz(ctx *gin.Context) {
	log.SetPrefix("http-server.handler.CreatePvz")
	var req request.CreatePvz
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if _, ok := entity.Cities[entity.City(req.City)]; !ok {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid city: "+req.City))
		return
	}

	pvz, err := h.pvzSrv.CreatePvz(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
	}
	ctx.JSON(http.StatusCreated, pvz.ToResponse())
}
