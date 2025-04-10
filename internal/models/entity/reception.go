package entity

type Status string

const (
	STATUS_IN_PROGRESS Status = "in_progress"
	STATUS_FINISHED    Status = "finished"
)

var Statuses map[Status]bool = map[Status]bool{
	STATUS_IN_PROGRESS: true,
	STATUS_FINISHED:    true,
}
