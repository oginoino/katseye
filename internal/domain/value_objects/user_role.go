package valueobjects

import "errors"

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleUser    UserRole = "user"
)

var (
	ErrInvalidUserRole = errors.New("invalid user role")
)

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) Validate() error {
	switch r {
	case UserRoleAdmin, UserRoleManager, UserRoleUser:
		return nil
	default:
		return ErrInvalidUserRole
	}
}
