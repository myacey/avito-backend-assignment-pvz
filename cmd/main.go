package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	http_server "github.com/myacey/avito-backend-assignment-pvz/internal/http-server"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
	"github.com/myacey/avito-backend-assignment-pvz/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	queries, conn, err := repository.ConfigurePostgres(cfg)

	receptionRepo := repository.NewReceptionRepository(queries)
	pvzRepo := repository.NewPvzRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	app := http_server.New(&service.Service{
		UserServiceImpl:      *service.NewUserService(*userRepo, conn),
		PvzServiceImpl:       *service.NewPvzService(*pvzRepo, conn),
		ReceptionServiceImpl: *service.NewReceptionService(*receptionRepo, conn),
	}, cfg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				app.Stop(ctx)
				return
			}
		}
	}()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
