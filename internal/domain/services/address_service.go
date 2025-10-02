package services

import (
	"context"
	"errors"
	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressService struct {
	addressRepo repositories.AddressRepository
}

func NewAddressService(addressRepo repositories.AddressRepository) *AddressService {
	return &AddressService{
		addressRepo: addressRepo,
	}
}

func (s *AddressService) GetAddressByID(ctx context.Context, id primitive.ObjectID) (*entities.Address, error) {
	return s.addressRepo.GetAddressByID(ctx, id)
}

func (s *AddressService) CreateAddress(ctx context.Context, address *entities.Address) error {
	if address == nil {
		return errors.New("address is nil")
	}

	// Validate address
	if err := address.Validate(); err != nil {
		return err
	}

	// Generate ID if not set
	if address.ID.IsZero() {
		address.ID = primitive.NewObjectID()
	}
	return s.addressRepo.CreateAddress(ctx, address)
}

func (s *AddressService) UpdateAddress(ctx context.Context, address *entities.Address) error {
	if address == nil {
		return errors.New("address is nil")
	}

	// Validate address
	if err := address.Validate(); err != nil {
		return err
	}

	return s.addressRepo.UpdateAddress(ctx, address)
}

func (s *AddressService) DeleteAddress(ctx context.Context, id primitive.ObjectID) error {
	return s.addressRepo.DeleteAddress(ctx, id)
}

func (s *AddressService) ListAddresses(ctx context.Context, filter map[string]interface{}) ([]*entities.Address, error) {
	return s.addressRepo.ListAddresses(ctx, filter)
}
