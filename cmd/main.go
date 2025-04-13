package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	pvzv1 "github.com/myacey/avito-backend-assignment-pvz/internal/grpc/pvz/v1"
	http_server "github.com/myacey/avito-backend-assignment-pvz/internal/http-server"
	"github.com/myacey/avito-backend-assignment-pvz/internal/http-server/handler"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/auth"
	jwt_token "github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwt-token"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
	"github.com/myacey/avito-backend-assignment-pvz/pkg/openapi"
	middleware "github.com/oapi-codegen/gin-middleware"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	queries, conn, err := repository.ConfigurePostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	receptionRepo := repository.NewReceptionRepository(queries)
	pvzRepo := repository.NewPvzRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	tokenSrv := jwt_token.New(cfg.TokenService)
	authSrv := auth.New(tokenSrv)

	pvzSrv := *service.NewPvzService(pvzRepo)
	srv := service.Service{
		UserService:      *service.NewUserService(userRepo, conn, tokenSrv),
		PvzService:       pvzSrv,
		ReceptionService: *service.NewReceptionService(receptionRepo, conn, &pvzSrv),
	}
	app := http_server.New(&srv, cfg, authSrv)

	hndlr := handler.NewHandler(&srv.ReceptionService, &srv.PvzService, &srv.UserService)

	swagger, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	openapi.RegisterHandlers(app.Router, hndlr)
	app.Router.Use(middleware.OapiRequestValidator(swagger))

	go func() {
		<-ctx.Done()
		app.Stop(ctx)
	}()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pvzv1.RegisterPVZServiceServer(grpcServer, pvzv1.NewPVZServerGRPC(&pvzSrv))
	go func() {
		log.Println("start grpc server on :3000")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
