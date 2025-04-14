package response

import (
	"time"

	"github.com/google/uuid"
)

type Error struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Code      int    `json:"code"`
}

type Login struct {
	Token string `json:"token"`
}

type Pvz struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"date_time"`
	ProductType string    `json:"type"`
	ReceptionID uuid.UUID `json:"reception_id"`
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzID    uuid.UUID `json:"pvz_id"`
	Status   string    `json:"status"`
}

type PvzWithReception struct {
	Pvz        *Pvz         `json:"pvz"`
	Receptions []*Reception `json:"receptions"`
}
