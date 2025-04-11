package http_server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
)

// RoleCheckerMiddleware is middleware interface
// for checking account's roles.
type RoleCheckerMiddleware interface {
	AuthMiddleware(neededRole entity.Role) gin.HandlerFunc
}

type App struct {
	server  web.Server
	router  *gin.Engine
	service *service.Service

	authService RoleCheckerMiddleware
}

func New(service *service.Service, cfg config.AppConfig, authService RoleCheckerMiddleware) *App {
	app := &App{
		service:     service,
		authService: authService,
	}
	app.initRoutes()
	app.server = web.NewServer(cfg.ServerCfg, app.router)

	return app
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) initRoutes() {
	app.router = gin.Default()

	app.router.POST("/dummyLogin", mappedHandler[handler.UserService](&app.service.UserService, handler.DummyLogin))
	app.router.POST("/login", mappedHandler[handler.UserService](&app.service.UserService, handler.Login))
	app.router.POST("/register", mappedHandler[handler.UserService](&app.service.UserService, handler.Register))

	employeeOnly := app.router.Group("/")
	employeeOnly.Use(app.authService.AuthMiddleware(entity.ROLE_EMPLOYEE))
	{
		employeeOnly.POST("/pvz/:pvzid/close_last_reception", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, handler.FinishReception))
		employeeOnly.POST("/pvz/:pvzid/delete_last_product", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, handler.DeleteLastProduct))
		employeeOnly.POST("/receptions", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, handler.CreateReception))
		employeeOnly.POST("/products", mappedHandler[handler.ReceptionService](&app.service.ReceptionService, handler.AddProductToReception))
	}

	moderatorOnly := app.router.Group("/")
	moderatorOnly.Use(app.authService.AuthMiddleware(entity.ROLE_MODERATOR))
	{
		moderatorOnly.POST("/pvz", mappedHandler[handler.PvzService](&app.service.PvzService, handler.CreatePvz))
	}

	// apiHandler := api_impl.NewAPIHandlerWithFuncs(app.service)
	// openapi.RegisterHandlers(app.router, apiHandler)
}

func mappedHandler[S any](service S, handler func(*gin.Context, S) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := handler(ctx, service); err != nil {
			if httpError, ok := err.(apperror.HTTPError); ok {
				ctx.JSON(httpError.Code, response.Error{
					Code:      httpError.Code,
					Message:   httpError.Error(),
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
			ctx.Abort()
		}
	}
}
