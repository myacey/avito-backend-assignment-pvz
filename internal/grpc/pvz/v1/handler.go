package pvzv1

import (
	context "context"
	"math"
	"time"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/request"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type PVZServer struct {
	UnimplementedPVZServiceServer
	srv PvzFinder
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
