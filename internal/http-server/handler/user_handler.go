package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
)

type UserService interface {
	DummyLogin(context.Context, *request.DummyLogin) (*response.Login, error)
	Register(context.Context, *request.Register) (*entity.User, error)
	Login(context.Context, *request.Login) (*response.Login, error)
}

func DummyLogin(ctx *gin.Context, service UserService) error {
	log.SetPrefix("http-server.handler.DummyLogin")

	var req request.DummyLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req: " + err.Error())
	}

	if _, ok := entity.Roles[entity.Role(req.Role)]; !ok {
		return apperror.NewBadReq("invalid role: " + req.Role)
	}

	resp, err := service.DummyLogin(ctx, &req)
	if err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, resp)
	return nil
}

func Register(ctx *gin.Context, service UserService) error {
	log.SetPrefix("http-server.handler.Register")

	var req request.Register
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req: " + err.Error())
	}

	if _, ok := entity.Roles[entity.Role(req.Role)]; !ok {
		return apperror.NewBadReq("invalid rol: " + req.Role)
	}

	usr, err := service.Register(ctx, &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, usr.ToResponse())
	return nil
}

func Login(ctx *gin.Context, service UserService) error {
	log.SetPrefix("http-server.handler.Login")

	var req request.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return apperror.NewBadReq("invalid req: " + err.Error())
	}

	resp, err := service.Login(ctx, &req)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}
