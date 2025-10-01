package entities

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidPassword indicates the provided password does not match the stored hash.
var ErrInvalidPassword = errors.New("invalid password")

// User represents an authenticated account within the system.
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Active       bool               `json:"active" bson:"active"`
}

// Normalize prepares user fields for persistence/lookup.
func (u *User) Normalize() {
	if u == nil {
		return
	}
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
}

// CheckPassword compares a clear-text password with the stored bcrypt hash.
func (u *User) CheckPassword(password string) error {
	if u == nil {
		return ErrInvalidPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

// IsActive returns true when the account is enabled for authentication.
func (u *User) IsActive() bool {
	if u == nil {
		return false
	}
	return u.Active
}
