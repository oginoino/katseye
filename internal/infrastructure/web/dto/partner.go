package dto

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"
)

// PartnerRequest modela o payload de entrada para parceiros.
type PartnerRequest struct {
	Name          string                     `json:"partner_name"`
	Type          string                     `json:"partner_type"`
	Attributes    entities.PartnerAttributes `json:"partner_attributes"`
	AcceptedTypes []string                   `json:"accepted_types"`
}

// PartnerResponse padroniza a saída HTTP de parceiros.
type PartnerResponse struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"partner_name"`
	Type          string                     `json:"partner_type"`
	Attributes    entities.PartnerAttributes `json:"partner_attributes"`
	AcceptedTypes []string                   `json:"accepted_types"`
}

// ToEntity converte o DTO em uma entidade de parceiro.
func (req *PartnerRequest) ToEntity(id primitive.ObjectID) (*entities.Partner, error) {
	if req == nil {
		return nil, fmt.Errorf("partner request is nil")
	}

	partnerType, err := valueobjects.NewPartnerType(req.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid partner type: %w", err)
	}

	acceptedTypes := make([]valueobjects.ProductType, 0, len(req.AcceptedTypes))
	for _, rawType := range req.AcceptedTypes {
		productType, err := valueobjects.NewProductType(rawType)
		if err != nil {
			return nil, fmt.Errorf("invalid accepted product type %q: %w", rawType, err)
		}
		acceptedTypes = append(acceptedTypes, productType)
	}

	partner := &entities.Partner{
		ID:            id,
		Name:          req.Name,
		Type:          partnerType,
		Attributes:    req.Attributes,
		AcceptedTypes: acceptedTypes,
	}

	return partner, nil
}

// NewPartnerResponse converte a entidade de domínio para DTO.
func NewPartnerResponse(partner *entities.Partner) PartnerResponse {
	if partner == nil {
		return PartnerResponse{}
	}

	accepted := make([]string, 0, len(partner.AcceptedTypes))
	for _, t := range partner.AcceptedTypes {
		accepted = append(accepted, t.String())
	}

	return PartnerResponse{
		ID:            partner.ID.Hex(),
		Name:          partner.Name,
		Type:          partner.Type.String(),
		Attributes:    partner.Attributes,
		AcceptedTypes: accepted,
	}
}

// NewPartnerResponseList cria uma lista de DTOs a partir de entidades.
func NewPartnerResponseList(partners []*entities.Partner) []PartnerResponse {
	if len(partners) == 0 {
		return nil
	}

	responses := make([]PartnerResponse, 0, len(partners))
	for _, partner := range partners {
		responses = append(responses, NewPartnerResponse(partner))
	}

	return responses
}
