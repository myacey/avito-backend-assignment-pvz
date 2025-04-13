package http_server

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
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

	middleware "github.com/oapi-codegen/gin-middleware"
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

	app.server = web.NewServer(cfg.ServerCfg, app.Router)

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

// func (app *App) initRoutes() {
// 	app.router = gin.Default()

// 	app.router.POST("/dummyLogin", mappedHandler[handler.UserService](&app.service.UserService, app.handler.DummyLogin))
// 	app.router.POST("/login", mappedHandler[handler.UserService](&app.service.UserService, app.handler.Login))
// 	app.router.POST("/register", mappedHandler[handler.UserService](&app.service.UserService, app.handler.Register))

// 	authorizedOnly := app.router.Group("/")
// 	authorizedOnly.Use(app.authService.AuthMiddleware(
// 		entity.ROLE_EMPLOYEE,
// 		entity.ROLE_MODERATOR,
// 	))
// 	{
// 		authorizedOnly.GET("/pvz", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, app.handler.SearchReceptions))
// 	}
// 	// app.router.GET("/pvz", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, app.handler.SearchReceptions))

// 	employeeOnly := app.router.Group("/")
// 	employeeOnly.Use(app.authService.AuthMiddleware(entity.ROLE_EMPLOYEE))
// 	{
// 		employeeOnly.POST("/pvz/:pvzid/close_last_reception", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, app.handler.FinishReception))
// 		employeeOnly.POST("/pvz/:pvzid/delete_last_product", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, app.handler.DeleteLastProduct))
// 		employeeOnly.POST("/receptions", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, app.handler.CreateReception))
// 		employeeOnly.POST("/products", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, app.handler.AddProductToReception))
// 	}

// 	moderatorOnly := app.router.Group("/")
// 	moderatorOnly.Use(app.authService.AuthMiddleware(entity.ROLE_MODERATOR))
// 	{
// 		moderatorOnly.POST("/pvz", mappedHandler[handler.PvzService](&app.service.PvzService, app.handler.CreatePvz))
// 	}

// 	// apiHandler := api_impl.NewAPIHandlerWithFuncs(app.service)
// 	// openapi.RegisterHandlers(app.router, apiHandler)
// }

// func mappedHandler[S any](service S, handler func(*gin.Context, S) error) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		if err := handler(ctx, service); err != nil {
// 			if httpError, ok := err.(apperror.HTTPError); ok {
// 				ctx.JSON(httpError.Code, response.Error{
// 					Code:      httpError.Code,
// 					Message:   httpError.Error(),
// 					RequestId: ctx.GetHeader("X-Request-Id"),
// 				})

// 				if httpError.Code == http.StatusInternalServerError {
// 					log.Printf("internal error: %v | %v", httpError.Message, httpError.DebugError)
// 				}
// 			} else {
// 				ctx.JSON(http.StatusInternalServerError, response.Error{
// 					Code:      http.StatusInternalServerError,
// 					Message:   err.Error(),
// 					RequestId: ctx.GetHeader("X-Request-Id"),
// 				})
// 			}
// 			ctx.Set("Retry-After", 10)
// 			ctx.Abort()
// 		}
// 	}
// }
