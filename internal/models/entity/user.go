package entity

import "database/sql/driver"

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
