package services

import (
	"context"
	"errors"
	"time"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrConsumerRepositoryUnavailable = errors.New("consumer repository unavailable")
	ErrConsumerNotFound              = errors.New("consumer not found")
	ErrProductRepositoryUnavailable  = errors.New("product repository unavailable")
	ErrProductNotFound               = errors.New("product not found")
	ErrConsumerUserAlreadyLinked     = errors.New("consumer already linked to user")
	ErrConsumerUserNotLinked         = errors.New("consumer user not linked")
)

type ConsumerService struct {
	consumerRepo repositories.ConsumerRepository
	productRepo  repositories.ProductRepository
}

func NewConsumerService(consumerRepo repositories.ConsumerRepository, productRepo repositories.ProductRepository) *ConsumerService {
	if consumerRepo == nil {
		return nil
	}

	return &ConsumerService{
		consumerRepo: consumerRepo,
		productRepo:  productRepo,
	}
}

func (s *ConsumerService) GetConsumerByID(ctx context.Context, id primitive.ObjectID) (*entities.Consumer, error) {
	if s == nil || s.consumerRepo == nil {
		return nil, ErrConsumerRepositoryUnavailable
	}
	if id.IsZero() {
		return nil, errors.New("consumer id is required")
	}

	return s.consumerRepo.GetConsumerByID(ctx, id)
}

func (s *ConsumerService) CreateConsumer(ctx context.Context, consumer *entities.Consumer) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if consumer == nil {
		return entities.ErrConsumerNil
	}

	if err := consumer.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()

	if consumer.ID.IsZero() {
		consumer.ID = primitive.NewObjectID()
	}

	if consumer.CreatedAt.IsZero() {
		consumer.CreatedAt = now
	}
	consumer.UpdatedAt = now

	return s.consumerRepo.CreateConsumer(ctx, consumer)
}

func (s *ConsumerService) UpdateConsumer(ctx context.Context, consumer *entities.Consumer) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if consumer == nil {
		return entities.ErrConsumerNil
	}
	if consumer.ID.IsZero() {
		return errors.New("consumer id is required")
	}

	if err := consumer.Validate(); err != nil {
		return err
	}

	consumer.UpdatedAt = time.Now().UTC()

	return s.consumerRepo.UpdateConsumer(ctx, consumer)
}

func (s *ConsumerService) DeleteConsumer(ctx context.Context, id primitive.ObjectID) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if id.IsZero() {
		return errors.New("consumer id is required")
	}

	return s.consumerRepo.DeleteConsumer(ctx, id)
}

func (s *ConsumerService) ListConsumers(ctx context.Context, filter map[string]interface{}) ([]*entities.Consumer, error) {
	if s == nil || s.consumerRepo == nil {
		return nil, ErrConsumerRepositoryUnavailable
	}

	return s.consumerRepo.ListConsumers(ctx, filter)
}

func (s *ConsumerService) ContractProduct(ctx context.Context, consumerID, productID primitive.ObjectID) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if consumerID.IsZero() {
		return errors.New("consumer id is required")
	}
	if productID.IsZero() {
		return errors.New("product id is required")
	}

	consumer, err := s.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		return err
	}
	if consumer == nil {
		return ErrConsumerNotFound
	}

	if s.productRepo == nil {
		return ErrProductRepositoryUnavailable
	}

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	if err := consumer.AddContractedProduct(productID); err != nil {
		return err
	}

	consumer.UpdatedAt = time.Now().UTC()

	return s.consumerRepo.UpdateConsumer(ctx, consumer)
}

func (s *ConsumerService) RemoveContractedProduct(ctx context.Context, consumerID, productID primitive.ObjectID) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if consumerID.IsZero() {
		return errors.New("consumer id is required")
	}
	if productID.IsZero() {
		return errors.New("product id is required")
	}

	consumer, err := s.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		return err
	}
	if consumer == nil {
		return ErrConsumerNotFound
	}

	if err := consumer.RemoveContractedProduct(productID); err != nil {
		return err
	}

	consumer.UpdatedAt = time.Now().UTC()

	return s.consumerRepo.UpdateConsumer(ctx, consumer)
}

func (s *ConsumerService) AttachUserProfile(ctx context.Context, consumerID, userID primitive.ObjectID) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if consumerID.IsZero() {
		return errors.New("consumer id is required")
	}
	if userID.IsZero() {
		return errors.New("user id is required")
	}

	consumer, err := s.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		return err
	}
	if consumer == nil {
		return ErrConsumerNotFound
	}

	if consumer.HasLinkedUser() {
		return ErrConsumerUserAlreadyLinked
	}

	consumer.UserID = userID
	consumer.UpdatedAt = time.Now().UTC()

	return s.consumerRepo.UpdateConsumer(ctx, consumer)
}

func (s *ConsumerService) DetachUserProfile(ctx context.Context, consumerID primitive.ObjectID) error {
	if s == nil || s.consumerRepo == nil {
		return ErrConsumerRepositoryUnavailable
	}
	if consumerID.IsZero() {
		return errors.New("consumer id is required")
	}

	consumer, err := s.consumerRepo.GetConsumerByID(ctx, consumerID)
	if err != nil {
		return err
	}
	if consumer == nil {
		return ErrConsumerNotFound
	}

	if !consumer.HasLinkedUser() {
		return ErrConsumerUserNotLinked
	}

	consumer.UserID = primitive.NilObjectID
	consumer.UpdatedAt = time.Now().UTC()

	return s.consumerRepo.UpdateConsumer(ctx, consumer)
}
