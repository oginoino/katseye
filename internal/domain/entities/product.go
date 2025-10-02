package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID
	Name        string
	Category    valueObjects.ProductCategory
	Attributes  ProductAttributes
	PartnerID   primitive.ObjectID
	ProductType valueObjects.ProductType
}

// Validate performs validation on the product entity
func (p *Product) Validate() error {
	if p == nil {
		return errors.New("product is nil")
	}

	if p.Name == "" {
		return errors.New("product name is required")
	}

	if p.PartnerID.IsZero() {
		return errors.New("partner id is required")
	}

	if err := p.ProductType.Validate(); err != nil {
		return err
	}

	if err := p.Attributes.Validate(p.ProductType); err != nil {
		return err
	}

	return nil
}
