package repositories

import (
	"context"
	"errors"

	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}
