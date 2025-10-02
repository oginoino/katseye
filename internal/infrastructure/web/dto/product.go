package dto

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"
)

// ProductRequest representa o payload de entrada para criação/atualização de produtos.
type ProductRequest struct {
	Name        string                     `json:"product_name"`
	Category    string                     `json:"product_category"`
	PartnerID   string                     `json:"partner_id"`
	ProductType string                     `json:"product_type"`
	Attributes  entities.ProductAttributes `json:"product_attributes"`
}

// ProductResponse padroniza a saída HTTP para recursos de produto.
type ProductResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"product_name"`
	Category    string                     `json:"product_category"`
	PartnerID   string                     `json:"partner_id"`
	ProductType string                     `json:"product_type"`
	Attributes  entities.ProductAttributes `json:"product_attributes"`
}

// ToEntity converte o DTO em uma entidade de domínio pronta para validação.
func (req *ProductRequest) ToEntity(id primitive.ObjectID) (*entities.Product, error) {
	if req == nil {
		return nil, fmt.Errorf("product request is nil")
	}

	partnerID, err := primitive.ObjectIDFromHex(req.PartnerID)
	if err != nil {
		return nil, fmt.Errorf("invalid partner id: %w", err)
	}

	productType, err := valueobjects.NewProductType(req.ProductType)
	if err != nil {
		return nil, fmt.Errorf("invalid product type: %w", err)
	}

	product := &entities.Product{
		ID:          id,
		Name:        req.Name,
		Category:    valueobjects.ProductCategory(req.Category),
		Attributes:  req.Attributes,
		PartnerID:   partnerID,
		ProductType: productType,
	}

	return product, nil
}

// NewProductResponse cria um DTO de resposta a partir da entidade de domínio.
func NewProductResponse(product *entities.Product) ProductResponse {
	if product == nil {
		return ProductResponse{}
	}

	response := ProductResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Category:    string(product.Category),
		PartnerID:   product.PartnerID.Hex(),
		ProductType: product.ProductType.String(),
		Attributes:  product.Attributes,
	}

	return response
}

// NewProductResponseList converte uma lista de entidades em DTOs de resposta.
func NewProductResponseList(products []*entities.Product) []ProductResponse {
	if len(products) == 0 {
		return nil
	}

	responses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, NewProductResponse(product))
	}

	return responses
}
