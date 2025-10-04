package dto

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"
)

// PartnerRequest modela o payload de entrada para parceiros.
type PartnerRequest struct {
	Name              string                     `json:"partner_name"`
	Type              string                     `json:"partner_type"`
	Attributes        entities.PartnerAttributes `json:"partner_attributes"`
	ManagerProfileIDs []string                   `json:"manager_profile_ids"`
	AcceptedTypes     []string                   `json:"accepted_types"`
}

// PartnerResponse padroniza a saída HTTP de parceiros.
type PartnerResponse struct {
	ID                string                     `json:"id"`
	Name              string                     `json:"partner_name"`
	Type              string                     `json:"partner_type"`
	Attributes        entities.PartnerAttributes `json:"partner_attributes"`
	ManagerProfileIDs []string                   `json:"manager_profile_ids"`
	AcceptedTypes     []string                   `json:"accepted_types"`
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

	managerProfileIDs := make([]primitive.ObjectID, 0, len(req.ManagerProfileIDs))
	for _, rawID := range req.ManagerProfileIDs {
		trimmed := strings.TrimSpace(rawID)
		if trimmed == "" {
			continue
		}
		managerID, err := primitive.ObjectIDFromHex(trimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid manager profile id %q: %w", rawID, err)
		}
		managerProfileIDs = append(managerProfileIDs, managerID)
	}

	if len(req.ManagerProfileIDs) > 0 && len(managerProfileIDs) == 0 {
		return nil, fmt.Errorf("manager_profile_ids must contain at least one valid identifier")
	}

	if len(managerProfileIDs) == 0 {
		return nil, fmt.Errorf("manager_profile_ids is required")
	}

	partner := &entities.Partner{
		ID:                id,
		Name:              req.Name,
		Type:              partnerType,
		Attributes:        req.Attributes,
		ManagerProfileIDs: managerProfileIDs,
		AcceptedTypes:     acceptedTypes,
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

	managerProfiles := make([]string, 0, len(partner.ManagerProfileIDs))
	for _, managerID := range partner.ManagerProfileIDs {
		managerProfiles = append(managerProfiles, managerID.Hex())
	}

	return PartnerResponse{
		ID:                partner.ID.Hex(),
		Name:              partner.Name,
		Type:              partner.Type.String(),
		Attributes:        partner.Attributes,
		ManagerProfileIDs: managerProfiles,
		AcceptedTypes:     accepted,
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
