package entity

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	STATUS_IN_PROGRESS Status = "in_progress"
	STATUS_FINISHED    Status = "finished"
)

var Statuses map[Status]bool = map[Status]bool{
	STATUS_IN_PROGRESS: true,
	STATUS_FINISHED:    true,
}

func (r Status) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *Status) Scan(value interface{}) error {
	*r = Status(string(value.([]byte)))
	return nil
}

type ProductType string

const (
	PRODUCT_TYPE_ELECTRONICS ProductType = "электроника"
	PRODUCT_TYPE_CLOTHES     ProductType = "одежда"
	PRODUCT_TYPE_SHOES       ProductType = "обувь"
)

var ProductTypes map[ProductType]bool = map[ProductType]bool{
	PRODUCT_TYPE_ELECTRONICS: true,
	PRODUCT_TYPE_CLOTHES:     true,
	PRODUCT_TYPE_SHOES:       true,
}

type Reception struct {
	ID       uuid.UUID
	DateTime time.Time
	PvzID    uuid.UUID
	Status   Status
}

type Product struct {
	ID          uuid.UUID
	DateTime    time.Time
	Type        ProductType
	ReceptionID uuid.UUID
}
