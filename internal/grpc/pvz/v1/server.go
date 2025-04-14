package pvzv1

import (
	context "context"
	"errors"
	"log"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/entity"
)

type Config struct {
	Address          string        `mapstructure:"listen"`
	EnableTLS        bool          `mapstructure:"enable_tls"`
	CertFile         string        `mapstructure:"cert_file"`
	KeyFile          string        `mapstructure:"key_file"`
	KeepAliveTime    time.Duration `mapstructure:"keep_alive_time"`
	KeepAliveTimeout time.Duration `mapstructure:"keep_alive_timeout"`
}

type PvzFinder interface {
	SearchPvz(ctx context.Context, req *request.SearchPvz) ([]*entity.Pvz, error)
}

type Server struct {
	cfg    Config
	srv    PvzFinder
	server *grpc.Server
	lis    net.Listener
}

func New(cfg Config, service PvzFinder) (*Server, error) {
	if service == nil {
		return nil, errors.New("pvz service can't be nil")
	}

	var options []grpc.ServerOption
	options = append(options, grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    cfg.KeepAliveTime,
		Timeout: cfg.KeepAliveTimeout,
	}))

	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}
		options = append(options, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(options...)
	handler := &PVZServer{srv: service}
	RegisterPVZServiceServer(grpcServer, handler)

	return &Server{
		cfg:    cfg,
		srv:    service,
		server: grpcServer,
	}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	s.lis = lis

	log.Printf("starting gRPC server on %s", s.cfg.Address)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("gRPC Serve error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() {
	log.Println("shutting down gRPC server...")
	s.server.GracefulStop()
}
