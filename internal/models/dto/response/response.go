package response

import (
	"time"

	"github.com/google/uuid"
)

type Error struct {
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
}

type Login struct {
	Token string `json:"token"`
}

type CreatePvz struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzId    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"`
}
