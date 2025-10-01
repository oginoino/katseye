package services

import (
	"context"
	"errors"
	"katseye/internal/application/interfaces/repositories"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PartnerService struct {
	partnerRepo repositories.PartnerRepository
}

func NewPartnerService(partnerRepo repositories.PartnerRepository) *PartnerService {
	return &PartnerService{
		partnerRepo: partnerRepo,
	}
}

func (s *PartnerService) GetPartnerByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error) {
	return s.partnerRepo.GetPartnerByID(ctx, id)
}

func (s *PartnerService) CreatePartner(ctx context.Context, partner *entities.Partner) error {
	if partner == nil {
		return errors.New("partner is nil")
	}

	// Validate partner
	if err := partner.Validate(); err != nil {
		return err
	}

	// Generate ID if not set
	if partner.ID.IsZero() {
		partner.ID = primitive.NewObjectID()
	}
	return s.partnerRepo.CreatePartner(ctx, partner)
}

func (s *PartnerService) UpdatePartner(ctx context.Context, partner *entities.Partner) error {
	if partner == nil {
		return errors.New("partner is nil")
	}

	// Validate partner
	if err := partner.Validate(); err != nil {
		return err
	}

	return s.partnerRepo.UpdatePartner(ctx, partner)
}

func (s *PartnerService) DeletePartner(ctx context.Context, id primitive.ObjectID) error {
	return s.partnerRepo.DeletePartner(ctx, id)
}

func (s *PartnerService) ListPartners(ctx context.Context, filter map[string]interface{}) ([]*entities.Partner, error) {
	return s.partnerRepo.ListPartners(ctx, filter)
}
