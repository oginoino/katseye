package services

import (
	"context"
	"errors"
	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"
	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	productRepo repositories.ProductRepository
	partnerRepo repositories.PartnerRepository
}

func NewProductService(productRepo repositories.ProductRepository, partnerRepo repositories.PartnerRepository) *ProductService {
	if productRepo == nil {
		return nil
	}

	return &ProductService{
		productRepo: productRepo,
		partnerRepo: partnerRepo,
	}
}

func (s *ProductService) GetProductByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	return s.productRepo.GetProductByID(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, product *entities.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	// Validate product
	if err := product.Validate(); err != nil {
		return err
	}

	if err := s.ensurePartnerAccepts(ctx, product.PartnerID, product.ProductType); err != nil {
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
		return errors.New("product is nil")
	}

	// Validate product
	if err := product.Validate(); err != nil {
		return err
	}

	if err := s.ensurePartnerAccepts(ctx, product.PartnerID, product.ProductType); err != nil {
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

var (
	ErrPartnerRepositoryUnavailable = errors.New("partner repository unavailable")
	ErrPartnerNotFound              = errors.New("partner not found")
	ErrProductTypeNotAccepted       = errors.New("product type not accepted by partner")
)

func (s *ProductService) ensurePartnerAccepts(ctx context.Context, partnerID primitive.ObjectID, productType valueobjects.ProductType) error {
	if partnerID.IsZero() {
		return errors.New("partner id is required")
	}
	if s == nil || s.partnerRepo == nil {
		return ErrPartnerRepositoryUnavailable
	}

	partner, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return err
	}
	if partner == nil {
		return ErrPartnerNotFound
	}

	for _, accepted := range partner.AcceptedTypes {
		if accepted == productType {
			return nil
		}
	}

	return ErrProductTypeNotAccepted
}
