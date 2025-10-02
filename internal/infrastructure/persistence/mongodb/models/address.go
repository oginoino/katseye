package models

import (
	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddressDocument representa o documento de endereço no MongoDB.
type AddressDocument struct {
	ID               primitive.ObjectID       `bson:"_id,omitempty"`
	Country          string                   `bson:"country"`
	State            string                   `bson:"state"`
	City             string                   `bson:"city"`
	District         string                   `bson:"district"`
	Street           string                   `bson:"street"`
	Number           string                   `bson:"number"`
	Complement       string                   `bson:"complement"`
	PostalCode       string                   `bson:"postal_code"`
	Type             valueobjects.AddressType `bson:"type"`
	FormattedAddress string                   `bson:"formatted_address"`
}

// ToEntity converte o documento em entidade de domínio.
func (doc AddressDocument) ToEntity() *entities.Address {
	return &entities.Address{
		ID:               doc.ID,
		Country:          doc.Country,
		State:            doc.State,
		City:             doc.City,
		District:         doc.District,
		Street:           doc.Street,
		Number:           doc.Number,
		Complement:       doc.Complement,
		PostalCode:       doc.PostalCode,
		Type:             doc.Type,
		FormattedAddress: doc.FormattedAddress,
	}
}

// NewAddressDocument cria o documento persistido a partir da entidade.
func NewAddressDocument(address *entities.Address) AddressDocument {
	if address == nil {
		return AddressDocument{}
	}

	return AddressDocument{
		ID:               address.ID,
		Country:          address.Country,
		State:            address.State,
		City:             address.City,
		District:         address.District,
		Street:           address.Street,
		Number:           address.Number,
		Complement:       address.Complement,
		PostalCode:       address.PostalCode,
		Type:             address.Type,
		FormattedAddress: address.FormattedAddress,
	}
}
