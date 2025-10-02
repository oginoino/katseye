package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Partner struct {
	ID            primitive.ObjectID
	Name          string
	Type          valueObjects.PartnerType
	Attributes    PartnerAttributes
	AcceptedTypes []valueObjects.ProductType
}

type PartnerData struct {
	PartnerID primitive.ObjectID
	Address   Address
}

// Validate performs validation on the partner entity
func (p *Partner) Validate() error {
	if p == nil {
		return errors.New("partner is nil")
	}

	if p.Name == "" {
		return errors.New("partner name is required")
	}

	if err := p.Type.Validate(); err != nil {
		return err
	}

	if err := p.Attributes.Validate(p.Type); err != nil {
		return err
	}

	// Validate accepted types
	for _, productType := range p.AcceptedTypes {
		if err := productType.Validate(); err != nil {
			return errors.New("invalid product type in accepted types: " + err.Error())
		}
	}

	return nil
}
