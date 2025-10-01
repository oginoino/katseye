package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                primitive.ObjectID           `json:"id" bson:"_id, unique"`
	ProductName       string                       `json:"product_name" bson:"product_name"`
	ProductCategory   valueObjects.ProductCategory `json:"product_category" bson:"product_category"`
	ProductAttributes ProductAttributes            `json:"product_attributes" bson:"product_attributes"`
	ProductPartner    Partner                      `json:"product_partner" bson:"product_partner"`
	ProductType       valueObjects.ProductType     `json:"product_type" bson:"product_type"`
}

// Validate performs validation on the product entity
func (p *Product) Validate() error {
	if p.ProductName == "" {
		return errors.New("product_name is required")
	}

	if err := p.ProductType.Validate(); err != nil {
		return err
	}

	if err := p.ProductAttributes.Validate(p.ProductType); err != nil {
		return err
	}

	// Validate that the product type is accepted by the partner
	partnerAccepts := false
	for _, acceptedType := range p.ProductPartner.AcceptedTypes {
		if acceptedType == p.ProductType {
			partnerAccepts = true
			break
		}
	}

	if !partnerAccepts {
		return errors.New("product type is not accepted by the partner")
	}

	return nil
}
