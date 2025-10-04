package entities

import (
	"errors"

	valueObjects "katseye/internal/domain/value_objects"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Partner struct {
	ID                primitive.ObjectID
	Name              string
	Type              valueObjects.PartnerType
	Attributes        PartnerAttributes
	ManagerProfileIDs []primitive.ObjectID
	AcceptedTypes     []valueObjects.ProductType
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

	if len(p.ManagerProfileIDs) == 0 {
		return errors.New("partner manager profiles are required")
	}

	seenManagers := make(map[primitive.ObjectID]struct{}, len(p.ManagerProfileIDs))
	for _, managerID := range p.ManagerProfileIDs {
		if managerID.IsZero() {
			return errors.New("partner manager profiles must contain valid user identifiers")
		}
		if _, exists := seenManagers[managerID]; exists {
			return errors.New("partner manager profiles must be unique")
		}
		seenManagers[managerID] = struct{}{}
	}

	// Validate accepted types
	for _, productType := range p.AcceptedTypes {
		if err := productType.Validate(); err != nil {
			return errors.New("invalid product type in accepted types: " + err.Error())
		}
	}

	return nil
}

// HasManagerProfile reports whether the partner already references the given manager user id.
func (p *Partner) HasManagerProfile(userID primitive.ObjectID) bool {
	if p == nil || userID.IsZero() {
		return false
	}

	for _, existing := range p.ManagerProfileIDs {
		if existing == userID {
			return true
		}
	}

	return false
}

// RemoveManagerProfile removes the given manager from the partner, returning true when it was present.
func (p *Partner) RemoveManagerProfile(userID primitive.ObjectID) bool {
	if p == nil || userID.IsZero() {
		return false
	}

	filtered := make([]primitive.ObjectID, 0, len(p.ManagerProfileIDs))
	removed := false
	for _, existing := range p.ManagerProfileIDs {
		if existing == userID {
			removed = true
			continue
		}
		filtered = append(filtered, existing)
	}

	if removed {
		p.ManagerProfileIDs = filtered
	}

	return removed
}
