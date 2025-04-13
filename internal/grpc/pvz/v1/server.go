package pvzv1

import (
	context "context"
	"math"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/timestamppb"

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

type PVZServer struct {
	UnimplementedPVZServiceServer
	srv PvzFinder
}

func New(cfg Config) (*grpc.Server, error) {
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
	return grpcServer, nil
}

func NewPVZServer(s PvzFinder) *PVZServer {
	return &PVZServer{srv: s}
}

func (s *PVZServer) GetPVZList(ctx context.Context, _ *GetPVZListRequest) (*GetPVZListResponse, error) {
	pvzs, err := s.srv.SearchPvz(ctx, &request.SearchPvz{
		StartDate: time.Date(0, 0, 0, 0, 0, 0, 0, time.Local),
		EndDate:   time.Now(),
		Page:      1,
		Limit:     math.MaxInt32,
	})
	if err != nil {
		return nil, err
	}

	var res []*PVZ
	for _, pvz := range pvzs {
		res = append(res, &PVZ{
			Id:               pvz.ID.String(),
			City:             string(pvz.City),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
		})
	}

	return &GetPVZListResponse{Pvzs: res}, nil
}
