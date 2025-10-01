package services

import (
	"context"
	"katseye/internal/application/interfaces/repositories"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) GetProductByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	return s.productRepo.GetProductByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *entities.Product) error {
	if product == nil {
		return nil
	}

	// Validate product
	if err := product.Validate(); err != nil {
		return err
	}

	// Generate ID if not set
	if product.ID.IsZero() {
		product.ID = primitive.NewObjectID()
	}
	return s.productRepo.CreateProduct(ctx, product)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *entities.Product) error {
	if product == nil {
		return nil
	}

	// Validate product
	if err := product.Validate(); err != nil {
		return err
	}

	return s.productRepo.UpdateProduct(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	return s.productRepo.DeleteProduct(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, filter map[string]interface{}) ([]*entities.Product, error) {
	return s.productRepo.ListProducts(ctx, filter)
}
