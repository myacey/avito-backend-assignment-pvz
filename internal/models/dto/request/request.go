package request

import (
	"time"

	"github.com/google/uuid"
)

type DummyLogin struct {
	Role string `json:"role" binding:"required"`
}

type Register struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreatePvz struct {
	ID               uuid.UUID `json:"id" binding:"required,uuid"`
	RegistrationDate time.Time `json:"registration_date" binding:"required"`
	City             string    `json:"city"`
}

type SearchPvz struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}

type CreateReception struct {
	PvzID uuid.UUID `json:"pvz_id" binding:"required,uuid"`
}

type AddProduct struct {
	Type  string    `json:"type" binding:"required"`
	PvzID uuid.UUID `json:"pvz_id" binding:"required,uuid"`
}
