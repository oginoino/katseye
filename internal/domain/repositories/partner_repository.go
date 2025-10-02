package repositories

import (
	"context"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PartnerRepository interface {
	GetPartnerByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error)
	CreatePartner(ctx context.Context, partner *entities.Partner) error
	UpdatePartner(ctx context.Context, partner *entities.Partner) error
	DeletePartner(ctx context.Context, id primitive.ObjectID) error
	ListPartners(ctx context.Context, filter map[string]interface{}) ([]*entities.Partner, error)
}
