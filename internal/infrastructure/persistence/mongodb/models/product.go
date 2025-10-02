package models

import (
	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductDocument representa o formato persistido de um produto na coleção do MongoDB.
type ProductDocument struct {
	ID            primitive.ObjectID           `bson:"_id,omitempty"`
	Name          string                       `bson:"product_name"`
	Category      valueobjects.ProductCategory `bson:"product_category"`
	Attributes    entities.ProductAttributes   `bson:"product_attributes"`
	PartnerID     primitive.ObjectID           `bson:"partner_id"`
	ProductType   valueobjects.ProductType     `bson:"product_type"`
	LegacyPartner *legacyPartnerDocument       `bson:"product_partner,omitempty"`
}

// ToEntity converte um documento do MongoDB em uma entidade de domínio.
func (doc ProductDocument) ToEntity() *entities.Product {
	partnerID := doc.PartnerID
	if partnerID.IsZero() && doc.LegacyPartner != nil {
		partnerID = doc.LegacyPartner.ID
	}

	return &entities.Product{
		ID:          doc.ID,
		Name:        doc.Name,
		Category:    doc.Category,
		Attributes:  doc.Attributes,
		PartnerID:   partnerID,
		ProductType: doc.ProductType,
	}
}

// NewProductDocument cria um documento a partir da entidade de domínio.
func NewProductDocument(product *entities.Product) ProductDocument {
	if product == nil {
		return ProductDocument{}
	}

	return ProductDocument{
		ID:          product.ID,
		Name:        product.Name,
		Category:    product.Category,
		Attributes:  product.Attributes,
		PartnerID:   product.PartnerID,
		ProductType: product.ProductType,
	}
}

type legacyPartnerDocument struct {
	ID            primitive.ObjectID         `bson:"_id,omitempty"`
	AcceptedTypes []valueobjects.ProductType `bson:"accepted_types"`
}
