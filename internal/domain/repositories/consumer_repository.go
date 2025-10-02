package repositories

import (
	"context"

	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsumerRepository interface {
	GetConsumerByID(ctx context.Context, id primitive.ObjectID) (*entities.Consumer, error)
	CreateConsumer(ctx context.Context, consumer *entities.Consumer) error
	UpdateConsumer(ctx context.Context, consumer *entities.Consumer) error
	DeleteConsumer(ctx context.Context, id primitive.ObjectID) error
	ListConsumers(ctx context.Context, filter map[string]interface{}) ([]*entities.Consumer, error)
}
