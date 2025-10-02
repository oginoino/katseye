package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID               primitive.ObjectID
	Country          string
	State            string
	City             string
	District         string
	Street           string
	Number           string
	Complement       string
	PostalCode       string
	Type             valueObjects.AddressType
	FormattedAddress string
}

// Validate performs validation on the address entity
func (a *Address) Validate() error {
	if a == nil {
		return errors.New("address is nil")
	}
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
		return errors.New("postal code is required")
	}
	if a.Type == "" {
		return errors.New("type is required")
	}
	// FormattedAddress can be optional or generated
	return nil
}
