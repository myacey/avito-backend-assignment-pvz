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
	"github.com/myacey/avito-backend-assignment-pvz/internal/repository"
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

	app := http_server.New(cfg, conn, queries)
	go func() {
		<-ctx.Done()
		app.Stop(ctx)
	}()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pvzv1.RegisterPVZServiceServer(grpcServer, pvzv1.NewPVZServerGRPC(&app.Service.PvzService))
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
