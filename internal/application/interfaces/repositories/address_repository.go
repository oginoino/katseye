package repositories

import (
	"context"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressRepository interface {
	GetAddressByID(ctx context.Context, id primitive.ObjectID) (*entities.Address, error)
	CreateAddress(ctx context.Context, address *entities.Address) error
	UpdateAddress(ctx context.Context, address *entities.Address) error
	DeleteAddress(ctx context.Context, id primitive.ObjectID) error
	ListAddresses(ctx context.Context, filter map[string]interface{}) ([]*entities.Address, error)
}