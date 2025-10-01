package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID               primitive.ObjectID       `json:"id" bson:"_id, unique"`
	Country          string                   `json:"country" bson:"country"`
	State            string                   `json:"state" bson:"state"`
	City             string                   `json:"city" bson:"city"`
	District         string                   `json:"district" bson:"district"`
	Street           string                   `json:"street" bson:"street"`
	Number           string                   `json:"number" bson:"number"`
	Complement       string                   `json:"complement" bson:"complement"`
	PostalCode       string                   `json:"postal_code" bson:"postal_code"`
	Type             valueObjects.AddressType `json:"type" bson:"type"`
	FormattedAddress string                   `json:"formatted_address" bson:"formatted_address"`
}

// Validate performs validation on the address entity
func (a *Address) Validate() error {
	if a.Country == "" {
		return errors.New("country is required")
	}
	if a.State == "" {
		return errors.New("state is required")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if a.Street == "" {
		return errors.New("street is required")
	}
	if a.Number == "" {
		return errors.New("number is required")
	}
	if a.PostalCode == "" {
		return errors.New("postal_code is required")
	}
	if a.Type == "" {
		return errors.New("type is required")
	}
	// FormattedAddress can be optional or generated
	return nil
}

