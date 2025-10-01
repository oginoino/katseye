package mongodb

import (
	"context"
	"errors"
	"katseye/internal/application/interfaces/repositories"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartnerRepositoryMongo struct {
	collection *mongo.Collection
}

func NewPartnerRepositoryMongo(collection *mongo.Collection) repositories.PartnerRepository {
	return &PartnerRepositoryMongo{
		collection: collection,
	}
}

func (r *PartnerRepositoryMongo) GetPartnerByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error) {
	var partner entities.Partner
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&partner)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Partner not found
		}
		return nil, err
	}
	return &partner, nil
}

func (r *PartnerRepositoryMongo) CreatePartner(ctx context.Context, partner *entities.Partner) error {
	_, err := r.collection.InsertOne(ctx, partner)
	return err
}

func (r *PartnerRepositoryMongo) UpdatePartner(ctx context.Context, partner *entities.Partner) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": partner.ID}, bson.M{"$set": partner})
	return err
}

func (r *PartnerRepositoryMongo) DeletePartner(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *PartnerRepositoryMongo) ListPartners(ctx context.Context, filter map[string]interface{}) ([]*entities.Partner, error) {
	var partners []*entities.Partner
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var partner entities.Partner
		if err := cursor.Decode(&partner); err != nil {
			return nil, err
		}
		partners = append(partners, &partner)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return partners, nil
}
