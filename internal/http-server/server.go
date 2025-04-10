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
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/auth"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/web/apperror"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
)

type App struct {
	server  web.Server
	router  *gin.Engine
	service *service.Service
}

func New(service *service.Service, cfg config.AppConfig) *App {
	app := &App{
		service: service,
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

	app.router.POST("/dummyLogin", mappedHandler[handler.UserService](app.service, handler.DummyLogin))
	app.router.POST("/login", mappedHandler[handler.UserService](app.service, handler.Login))
	app.router.POST("/register", mappedHandler[handler.UserService](app.service, handler.Login))

	employeeOnly := app.router.Group("/")
	employeeOnly.Use(auth.AuthMiddleware(entity.ROLE_EMPLOYEE))
	{
		employeeOnly.POST("/pvz/:pvzid/close_last_reception", mappedHandler[handler.ReceptionService](app.service, handler.CompleteReception))
		employeeOnly.POST("/pvz/:pvzid/delete_last_product", mappedHandler[handler.ReceptionService](app.service, handler.DeleteLastProduct))
		employeeOnly.POST("/receptions", mappedHandler[handler.ReceptionService](app.service, handler.CreateReception))
		employeeOnly.POST("/products", mappedHandler[handler.ReceptionService](app.service, handler.AddProductToReception))
	}

	moderatorOnly := app.router.Group("/")
	moderatorOnly.Use(auth.AuthMiddleware(entity.ROLE_MODERATOR))
	{
		moderatorOnly.POST("/pvz", mappedHandler[handler.PvzService](app.service, handler.CreatePvz))
	}
}

func mappedHandler[S any](service S, handler func(*gin.Context, S) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := handler(ctx, service); err != nil {
			if httpError, ok := err.(*apperror.HTTPError); ok {
				ctx.JSON(httpError.Code, response.Error{
					Code:      httpError.Code,
					Message:   httpError.Error(),
					RequestId: ctx.GetHeader("X-Request-Id"),
				})

				if httpError.Code == http.StatusInternalServerError {
					log.Printf("internal error: %v", httpError.Message)
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
