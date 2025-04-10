package entity

import (
	"time"

	"github.com/google/uuid"
)

type City string

const (
	CITY_MOSCOW           City = "Москва"
	CITY_SAINT_PETERSBURG City = "Санкт-Петербург"
	CITY_KAZAN            City = "Казань"
)

var Cities map[City]bool = map[City]bool{
	CITY_MOSCOW:           true,
	CITY_SAINT_PETERSBURG: true,
	CITY_KAZAN:            true,
}

type Pvz struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string
}
