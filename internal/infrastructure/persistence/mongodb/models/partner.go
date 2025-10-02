package models

import (
	"katseye/internal/domain/entities"
	valueobjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PartnerDocument descreve a estrutura do documento de parceiro no MongoDB.
type PartnerDocument struct {
	ID            primitive.ObjectID         `bson:"_id,omitempty"`
	Name          string                     `bson:"partner_name"`
	Type          string                     `bson:"partner_type"`
	Attributes    entities.PartnerAttributes `bson:"partner_attributes"`
	AcceptedTypes []string                   `bson:"accepted_types"`
}

// ToEntity converte o documento em uma entidade de domínio.
func (doc PartnerDocument) ToEntity() *entities.Partner {
	var partnerType valueobjects.PartnerType
	if t, err := valueobjects.NewPartnerType(doc.Type); err == nil {
		partnerType = t
	}

	accepted := make([]valueobjects.ProductType, 0, len(doc.AcceptedTypes))
	for _, raw := range doc.AcceptedTypes {
		if productType, err := valueobjects.NewProductType(raw); err == nil {
			accepted = append(accepted, productType)
		}
	}

	return &entities.Partner{
		ID:            doc.ID,
		Name:          doc.Name,
		Type:          partnerType,
		Attributes:    doc.Attributes,
		AcceptedTypes: accepted,
	}
}

// NewPartnerDocument converte uma entidade de domínio para o modelo persistido.
func NewPartnerDocument(partner *entities.Partner) PartnerDocument {
	if partner == nil {
		return PartnerDocument{}
	}

	accepted := make([]string, 0, len(partner.AcceptedTypes))
	for _, productType := range partner.AcceptedTypes {
		accepted = append(accepted, productType.String())
	}

	return PartnerDocument{
		ID:            partner.ID,
		Name:          partner.Name,
		Type:          partner.Type.String(),
		Attributes:    partner.Attributes,
		AcceptedTypes: accepted,
	}
}
