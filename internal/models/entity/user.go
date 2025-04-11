package entity

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type Role string

const (
	ROLE_EMPLOYEE  Role = "employee"
	ROLE_MODERATOR Role = "moderator"
)

// for fast-checking
var Roles map[Role]bool = map[Role]bool{
	ROLE_EMPLOYEE:  true,
	ROLE_MODERATOR: true,
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *Role) Scan(value interface{}) error {
	*r = Role(string(value.([]byte)))
	return nil
}

type User struct {
	ID       uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     Role      `json:"role"`
}
