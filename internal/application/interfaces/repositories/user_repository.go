package repositories

import (
	"context"

	"katseye/internal/domain/entities"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}
