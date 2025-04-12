package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
)

type ReceptionService interface {
	SearchReceptions(context.Context, *request.SearchPvz) ([]*entity.PvzWithReception, error)
	FinishReception(context.Context, uuid.UUID) (*entity.Reception, error)
	DeleteLastProduct(context.Context, uuid.UUID) error
	CreateReception(context.Context, *request.CreateReception) (*entity.Reception, error)
	AddProductToReception(context.Context, *request.AddProduct) (*entity.Product, error)
}

func SearchReceptions(ctx *gin.Context, sevice ReceptionService) error {
	log.SetPrefix("http-server.handler.SearchPvz")

	startDateStr := ctx.Query("startDate")
	if startDateStr == "" {
		return apperror.NewBadReq("start date can't be nil")
	}
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return apperror.NewBadReq("invalid start date: " + startDateStr)
	}

	endDateStr := ctx.Query("endDate")
	if endDateStr == "" {
		return apperror.NewBadReq("end date can't be nil")
	}
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return apperror.NewBadReq("invalid end date: " + endDateStr)
	}

	page := 1
	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return apperror.NewBadReq("invalid page: " + pageStr)
		}
	}

	limit := 10
	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return apperror.NewBadReq("invalid limit: " + limitStr)
		}
	}

	req := &request.SearchPvz{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	pvzWithReceptions, err := sevice.SearchReceptions(ctx, req) // TODO: gen responses
	if err != nil {
		return err
	}
	resp := make([]*response.PvzWithReception, len(pvzWithReceptions))
	for i, v := range pvzWithReceptions {
		resp[i] = v.ToResponse()
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}

func FinishReception(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.CloseLastReception")

	pvzID, err := uuid.Parse(ctx.Param("pvzid"))
	if err != nil {
		return apperror.NewBadReq("invalid pvz id")
	}

	reception, err := service.FinishReception(ctx, pvzID)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, reception.ToResponse())
	return nil
}

func DeleteLastProduct(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.DeleteLastProduct")

	pvzID, err := uuid.Parse(ctx.Param("pvzid"))
	if err != nil {
		return apperror.NewBadReq("invalid pvz id")
	}

	err = service.DeleteLastProduct(ctx, pvzID)
	if err != nil {
		return err
	}

	ctx.Status(http.StatusOK)
	return nil
}

func CreateReception(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.CreateReception")

	var req request.CreateReception
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req: " + err.Error())
	}

	reception, err := service.CreateReception(ctx, &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, reception.ToResponse())
	return nil
}

func AddProductToReception(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.AddProductToReception")

	var req request.AddProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req: " + err.Error())
	}

	if _, ok := entity.ProductTypes[entity.ProductType(req.Type)]; !ok {
		return apperror.NewBadReq("invalid product type: " + req.Type)
	}

	product, err := service.AddProductToReception(ctx, &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, product.ToResponse())
	return nil
}
