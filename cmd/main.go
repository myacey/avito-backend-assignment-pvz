package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/myacey/avito-backend-assignment-pvz/internal/config"
	pvzv1 "github.com/myacey/avito-backend-assignment-pvz/internal/grpc/pvz/v1"
	"github.com/myacey/avito-backend-assignment-pvz/internal/httpserver"
	"github.com/myacey/avito-backend-assignment-pvz/internal/pkg/metrics"
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
)

var (
	cfgPath = flag.String("f", "configs/config.yaml", "path to the go's config")
	useGrpc = flag.Bool("g", true, "use grpc server?")
)

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	queries, conn, err := repository.ConfigurePostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := httpserver.New(cfg, conn, queries)

	var grpcServer *grpc.Server
	if *useGrpc {
		grpcServer, err := pvzv1.New(cfg.GRPCServerCfg, &app.Service.PvzService)
		if err != nil {
			log.Fatalf("failed to create grpc server: %v", err)
		}
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}

	go func() {
		<-ctx.Done()

		app.Stop(ctx)
		if grpcServer != nil {
			grpcServer.Stop()
		}
	}()

	metrics.StartMetricsServer()

	if err := app.Start(ctx); err != nil {
		log.Printf("http server returned error: %v", err)
	}
}
