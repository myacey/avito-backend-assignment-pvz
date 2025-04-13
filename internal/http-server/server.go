package http_server

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/auth"
	jwt_token "github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwt-token"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/metrics"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	db "github.com/myacey/avito-backend-assignment-pvz/internal/repository/sqlc"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
	"github.com/myacey/avito-backend-assignment-pvz/pkg/openapi"
)

type App struct {
	server  web.Server
	Router  *gin.Engine
	Service *service.Service
}

func New(cfg config.AppConfig, conn *sql.DB, queries *db.Queries) *App {
	app := &App{
		Router: gin.Default(),
	}
	app.initialize(cfg, conn, queries)

	app.server = web.NewServer(cfg.HTTPServerCfg, app.Router)

	return app
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) initialize(cfg config.AppConfig, conn *sql.DB, queries *db.Queries) {
	receptionRepo := repository.NewReceptionRepository(queries)
	pvzRepo := repository.NewPvzRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	tokenSrv := jwt_token.New(cfg.TokenService)
	authSrv := auth.New(tokenSrv)

	pvzSrv := *service.NewPvzService(pvzRepo)
	app.Service = &service.Service{
		UserService:      *service.NewUserService(userRepo, conn, tokenSrv),
		PvzService:       pvzSrv,
		ReceptionService: *service.NewReceptionService(receptionRepo, conn, &pvzSrv),
	}

	hndlr := handler.NewHandler(
		&app.Service.ReceptionService,
		&app.Service.PvzService,
		&app.Service.UserService,
		authSrv,
	)

	app.Router.Use(metrics.GetMetricsMiddleware())

	swagger, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}

	openapi.RegisterHandlers(app.Router, hndlr)
	app.Router.Use(middleware.OapiRequestValidator(swagger))
}
