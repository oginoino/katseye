package dto

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"
)

// AddressRequest encapsula o payload HTTP para o recurso de endereço.
type AddressRequest struct {
	Country          string `json:"country"`
	State            string `json:"state"`
	City             string `json:"city"`
	District         string `json:"district"`
	Street           string `json:"street"`
	Number           string `json:"number"`
	Complement       string `json:"complement"`
	PostalCode       string `json:"postal_code"`
	Type             string `json:"type"`
	FormattedAddress string `json:"formatted_address"`
}

// AddressResponse estrutura o payload de saída para endereços.
type AddressResponse struct {
	ID               string `json:"id"`
	Country          string `json:"country"`
	State            string `json:"state"`
	City             string `json:"city"`
	District         string `json:"district"`
	Street           string `json:"street"`
	Number           string `json:"number"`
	Complement       string `json:"complement"`
	PostalCode       string `json:"postal_code"`
	Type             string `json:"type"`
	FormattedAddress string `json:"formatted_address"`
}

// ToEntity converte o DTO em uma entidade de domínio.
func (req *AddressRequest) ToEntity(id primitive.ObjectID) (*entities.Address, error) {
	if req == nil {
		return nil, fmt.Errorf("address request is nil")
	}

	address := &entities.Address{
		ID:               id,
		Country:          req.Country,
		State:            req.State,
		City:             req.City,
		District:         req.District,
		Street:           req.Street,
		Number:           req.Number,
		Complement:       req.Complement,
		PostalCode:       req.PostalCode,
		Type:             valueobjects.AddressType(req.Type),
		FormattedAddress: req.FormattedAddress,
	}

	return address, nil
}

// NewAddressResponse cria um DTO para saída.
func NewAddressResponse(address *entities.Address) AddressResponse {
	if address == nil {
		return AddressResponse{}
	}

	return AddressResponse{
		ID:               address.ID.Hex(),
		Country:          address.Country,
		State:            address.State,
		City:             address.City,
		District:         address.District,
		Street:           address.Street,
		Number:           address.Number,
		Complement:       address.Complement,
		PostalCode:       address.PostalCode,
		Type:             string(address.Type),
		FormattedAddress: address.FormattedAddress,
	}
}

// NewAddressResponseList converte uma slice de entidades em DTOs.
func NewAddressResponseList(addresses []*entities.Address) []AddressResponse {
	if len(addresses) == 0 {
		return nil
	}

	responses := make([]AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		responses = append(responses, NewAddressResponse(address))
	}

	return responses
}
