package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	http_server "github.com/myacey/avito-backend-assignment-pvz/internal/http-server"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/auth"
	jwt_token "github.com/myacey/avito-backend-assignment-pvz/internal/pkg/jwt-token"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"

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

	app := http_server.New(&service.Service{
		UserService:      *service.NewUserService(*userRepo, conn, tokenSrv),
		PvzService:       *service.NewPvzService(*pvzRepo, conn),
		ReceptionService: *service.NewReceptionService(*receptionRepo, conn),
	}, cfg, authSrv)

	go func() {
		<-ctx.Done()
		app.Stop(ctx)
	}()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
