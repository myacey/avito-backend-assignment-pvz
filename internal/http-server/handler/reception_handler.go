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
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
)

type ReceptionService interface {
	SearchReceptions(context.Context, *request.SearchPvz) (map[string]interface{}, error)
	CompleteReception(context.Context, string) (*response.Reception, error)
	DeleteLastProduct(context.Context, string) error
	CreateReception(context.Context, *request.CreateReception) (*response.Reception, error)
	AddProductToReception(context.Context, *request.AddProduct) error
}

func SearchReceptions(ctx *gin.Context, sevice ReceptionService) error {
	log.SetPrefix("http-server.handler.SearchPvz")

	var query map[string]string
	if err := ctx.BindQuery(&query); err != nil {
		return apperror.NewBadReq("invalid req")
	}

	page, limit := 1, 10

	startDate, err := time.Parse(time.RFC3339, query["startDate"])
	if err != nil {
		return apperror.NewBadReq("invalid start date query param")
	}

	endDate, err := time.Parse(time.RFC3339, query["endDate"])
	if err != nil {
		return apperror.NewBadReq("invalid end date query param")
	}

	if p, ok := query["page"]; ok {
		page, err = strconv.Atoi(p)
		if err != nil {
			return apperror.NewBadReq("invalid page query param")
		}
	}

	if l, ok := query["limit"]; ok {
		limit, err = strconv.Atoi(l)
		if err != nil {
			return apperror.NewBadReq("invalid limit query param")
		}
	}

	req := &request.SearchPvz{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}

	resp, err := sevice.SearchReceptions(ctx, req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}

func CompleteReception(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.CloseLastReception")

	pvzID, err := uuid.Parse(ctx.Param("pvzid"))
	if err != nil {
		return apperror.NewBadReq("invalid pvz id")
	}

	resp, err := service.CompleteReception(ctx, pvzID.String())
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}

func DeleteLastProduct(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.DeleteLastProduct")

	pvzID, err := uuid.Parse(ctx.Param("pvzid"))
	if err != nil {
		return apperror.NewBadReq("invalid pvz id")
	}

	err = service.DeleteLastProduct(ctx, pvzID.String())
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
		return apperror.NewBadReq("invalid req")
	}

	resp, err := service.CreateReception(ctx, &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}

func AddProductToReception(ctx *gin.Context, service ReceptionService) error {
	log.SetPrefix("http-server.handler.AddProductToReception")

	var req request.AddProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req")
	}

	err := service.AddProductToReception(ctx, &req)
	if err != nil {
		return err
	}

	ctx.Status(http.StatusOK)
	return nil
}
