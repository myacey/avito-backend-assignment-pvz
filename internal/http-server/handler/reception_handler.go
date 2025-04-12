package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/pkg/openapi"
)

type ReceptionService interface {
	SearchReceptions(context.Context, *request.SearchPvz) ([]*entity.PvzWithReception, error)
	FinishReception(context.Context, uuid.UUID) (*entity.Reception, error)
	DeleteLastProduct(context.Context, uuid.UUID) error
	CreateReception(context.Context, *request.CreateReception) (*entity.Reception, error)
	AddProductToReception(context.Context, *request.AddProduct) (*entity.Product, error)
}

func (h Handler) GetPvz(ctx *gin.Context, params openapi.GetPvzParams) {
	log.SetPrefix("http-server.handler.SearchPvz")

	startDate := *params.StartDate
	endDate := *params.EndDate

	page, limit := 1, 10
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	req := &request.SearchPvz{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	pvzWithReceptions, err := h.receptionSrv.SearchReceptions(ctx, req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}
	resp := make([]*response.PvzWithReception, len(pvzWithReceptions))
	for i, v := range pvzWithReceptions {
		resp[i] = v.ToResponse()
	}

	ctx.JSON(http.StatusOK, resp)
}

// PostPvzPvzIdCloseLastReception ends reception
func (h Handler) PostPvzPvzIdCloseLastReception(ctx *gin.Context, pvzID uuid.UUID) {
	log.SetPrefix("http-server.handler.CloseLastReception")

	reception, err := h.receptionSrv.FinishReception(ctx, pvzID)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, reception.ToResponse())
}

// PostPvzPvzIdDeleteLastProduct deletes last product in reception
func (h Handler) PostPvzPvzIdDeleteLastProduct(ctx *gin.Context, pvzID uuid.UUID) {
	log.SetPrefix("http-server.handler.DeleteLastProduct")

	err := h.receptionSrv.DeleteLastProduct(ctx, pvzID)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// PostReceptiosn creates new reception on pvz.
func (h Handler) PostReceptions(ctx *gin.Context) {
	log.SetPrefix("http-server.handler.CreateReception")

	var req request.CreateReception
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	reception, err := h.receptionSrv.CreateReception(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, reception.ToResponse())
}

func (h Handler) PostProducts(ctx *gin.Context) {
	log.SetPrefix("http-server.handler.PostProducts")

	var req request.AddProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid req: "+err.Error()))
		return
	}

	if _, ok := entity.ProductTypes[entity.ProductType(req.Type)]; !ok {
		wrapCtxWithError(ctx, apperror.NewBadReq("invalid product type: "+req.Type))
		return
	}

	product, err := h.receptionSrv.AddProductToReception(ctx, &req)
	if err != nil {
		wrapCtxWithError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, product.ToResponse())
}
