package entities

import (
	"errors"
	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Partner struct {
	ID                primitive.ObjectID         `json:"id" bson:"_id, unique"`
	PartnerName       string                     `json:"partner_name" bson:"partner_name"`
	PartnerType       valueObjects.PartnerType   `json:"partner_type" bson:"partner_type"`
	PartnerAttributes PartnerAttributes          `json:"partner_attributes" bson:"partner_attributes"`
	AcceptedTypes     []valueObjects.ProductType `json:"accepted_types" bson:"accepted_types"`
}

type PartnerData struct {
	PartnerID primitive.ObjectID `json:"partner_id" bson:"partner_id"`
	Address   Address            `json:"address" bson:"address"`
}

// Validate performs validation on the partner entity
func (p *Partner) Validate() error {
	if p.PartnerName == "" {
		return errors.New("partner_name is required")
	}

	if err := p.PartnerType.Validate(); err != nil {
		return err
	}

	if err := p.PartnerAttributes.Validate(p.PartnerType); err != nil {
		return err
	}

	// Validate accepted types
	for _, productType := range p.AcceptedTypes {
		if err := productType.Validate(); err != nil {
			return errors.New("invalid product type in accepted_types: " + err.Error())
		}
	}

	return nil
}
