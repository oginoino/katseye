package repositories

import (
	"context"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error)
	CreateProduct(ctx context.Context, product *entities.Product) error
	UpdateProduct(ctx context.Context, product *entities.Product) error
	DeleteProduct(ctx context.Context, id primitive.ObjectID) error
	ListProducts(ctx context.Context, filter map[string]interface{}) ([]*entities.Product, error)
}
