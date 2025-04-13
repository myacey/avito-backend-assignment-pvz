package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	pvzv1 "github.com/myacey/avito-backend-assignment-pvz/internal/grpc/pvz/v1"
	http_server "github.com/myacey/avito-backend-assignment-pvz/internal/httpserver"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/metrics"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
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

	app := http_server.New(cfg, conn, queries)

	// grpcServer := grpc.NewServer()
	grpcServer, err := pvzv1.New(cfg.GRPCServerCfg)
	if err != nil {
		log.Fatalf("failed to create grpc server: %v", err)
	}
	pvzv1.RegisterPVZServiceServer(grpcServer, pvzv1.NewPVZServer(&app.Service.PvzService))

	lis, err := net.Listen("tcp", cfg.GRPCServerCfg.Address)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Printf("start grpc server on %s", cfg.GRPCServerCfg.Address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		<-ctx.Done()
		app.Stop(ctx)
		grpcServer.GracefulStop()
	}()

	metrics.StartMetricsServer()

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
