package services

import (
	"context"
	"errors"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrPartnerManagerAlreadyLinked = errors.New("partner manager already linked")
	ErrPartnerManagerNotLinked     = errors.New("partner manager not linked")
	ErrPartnerManagerRequired      = errors.New("partner must retain at least one manager profile")
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

func (s *PartnerService) AssignManagerProfile(ctx context.Context, partnerID, userID primitive.ObjectID) error {
	if s == nil || s.partnerRepo == nil {
		return ErrPartnerRepositoryUnavailable
	}
	if partnerID.IsZero() {
		return errors.New("partner id is required")
	}
	if userID.IsZero() {
		return errors.New("user id is required")
	}

	partner, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return err
	}
	if partner == nil {
		return ErrPartnerNotFound
	}

	if partner.HasManagerProfile(userID) {
		return ErrPartnerManagerAlreadyLinked
	}

	partner.ManagerProfileIDs = append(partner.ManagerProfileIDs, userID)

	if err := partner.Validate(); err != nil {
		return err
	}

	return s.partnerRepo.UpdatePartner(ctx, partner)
}

func (s *PartnerService) RemoveManagerProfile(ctx context.Context, partnerID, userID primitive.ObjectID) error {
	if s == nil || s.partnerRepo == nil {
		return ErrPartnerRepositoryUnavailable
	}
	if partnerID.IsZero() {
		return errors.New("partner id is required")
	}
	if userID.IsZero() {
		return errors.New("user id is required")
	}

	partner, err := s.partnerRepo.GetPartnerByID(ctx, partnerID)
	if err != nil {
		return err
	}
	if partner == nil {
		return ErrPartnerNotFound
	}

	if !partner.HasManagerProfile(userID) {
		return ErrPartnerManagerNotLinked
	}

	if len(partner.ManagerProfileIDs) <= 1 {
		return ErrPartnerManagerRequired
	}

	partner.RemoveManagerProfile(userID)

	if err := partner.Validate(); err != nil {
		return err
	}

	return s.partnerRepo.UpdatePartner(ctx, partner)
}
