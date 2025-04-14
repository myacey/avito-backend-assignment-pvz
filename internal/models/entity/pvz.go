package entity

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
)

type City string

const (
	CityMoscow          City = "Москва"
	CitySaintPetersburg City = "Санкт-Петербург"
	CityKazan           City = "Казань"
)

func (c *City) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*c = City(s)
	case string:
		*c = City(s)
	default:
		return fmt.Errorf("unsupported scan type for City: %v", src)
	}
	return nil
}

func (c City) Value() (driver.Value, error) {
	return string(c), nil
}

var Cities = map[City]bool{
	CityMoscow:          true,
	CitySaintPetersburg: true,
	CityKazan:           true,
}

type Pvz struct {
	ID               uuid.UUID
	RegistrationDate time.Time
	City             City
}

func (pvz *Pvz) ToResponse() *response.Pvz {
	return &response.Pvz{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             string(pvz.City),
	}
}

func (pvz *Pvz) MarshalJSON() ([]byte, error) {
	return nil, errors.New("entity.Pvz: direct JSON serialization forbidden, use response.Pvz")
}

type PvzWithReception struct {
	Pvz        *Pvz
	Receptions []*Reception
}

func (pvzwr *PvzWithReception) ToResponse() *response.PvzWithReception {
	r := make([]*response.Reception, len(pvzwr.Receptions))
	for i, v := range pvzwr.Receptions {
		r[i] = v.ToResponse()
	}

	return &response.PvzWithReception{
		Pvz:        pvzwr.Pvz.ToResponse(),
		Receptions: r,
	}
}

func (pvzwr *PvzWithReception) MarshalJSON() ([]byte, error) {
	return nil, errors.New("entity.PvzWithReception: direct JSON serialization forbidden, use recponse.PvzWithReception")
}
