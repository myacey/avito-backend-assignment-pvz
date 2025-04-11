package entity

import (
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

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"date_time"`
	PvzID    uuid.UUID `json:"pvz_id"`
	Status   Status    `json:"status"`
}
