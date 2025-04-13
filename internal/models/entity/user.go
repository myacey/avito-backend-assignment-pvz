package entity

import (
	"database/sql/driver"
	"errors"

	"github.com/google/uuid"

	"github.com/myacey/avito-backend-assignment-pvz/internal/models/dto/response"
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
	ID       uuid.UUID
	Email    string
	Password string
	Role     Role
}

func (u *User) ToResponse() *response.User {
	return &response.User{
		ID:    u.ID,
		Email: u.Email,
		Role:  string(u.Role),
	}
}

func (u *User) MarshalJSON() ([]byte, error) {
	return nil, errors.New("entity.User: direct JSON serialization forbidden, use response.User")
}
