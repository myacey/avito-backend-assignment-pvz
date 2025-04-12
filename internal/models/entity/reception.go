package entity

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
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

func (c *ProductType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*c = ProductType(s)
	case string:
		*c = ProductType(s)
	default:
		return fmt.Errorf("unsupported scan type for ProductType: %T", src)
	}
	return nil
}

func (c ProductType) Value() (driver.Value, error) {
	return string(c), nil
}

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

func (r *Reception) ToResponse() *response.Reception {
	return &response.Reception{
		ID:       r.ID,
		DateTime: r.DateTime,
		PvzID:    r.PvzID,
		Status:   string(r.Status),
	}
}

func (r *Reception) MarshalJSON() ([]byte, error) {
	return nil, errors.New("entity.Reception: direct JSON serialization forbidden, use response.Reception")
}

type Product struct {
	ID          uuid.UUID
	DateTime    time.Time
	Type        ProductType
	ReceptionID uuid.UUID
}

func (p *Product) ToResponse() *response.Product {
	return &response.Product{
		ID:          p.ID,
		DateTime:    p.DateTime,
		ProductType: string(p.Type),
		ReceptionID: p.ReceptionID,
	}
}

func (p *Product) MarshalJSON() ([]byte, error) {
	return nil, errors.New("entity.Product: direct JSON serialization forbidden, use response.Product")
}
