package mongodb

import (
	"context"
	"errors"
	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"
	"katseye/internal/infrastructure/persistence/mongodb/models"

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
	var doc models.PartnerDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Partner not found
		}
		return nil, err
	}
	return doc.ToEntity(), nil
}

func (r *PartnerRepositoryMongo) CreatePartner(ctx context.Context, partner *entities.Partner) error {
	_, err := r.collection.InsertOne(ctx, models.NewPartnerDocument(partner))
	return err
}

func (r *PartnerRepositoryMongo) UpdatePartner(ctx context.Context, partner *entities.Partner) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": partner.ID}, bson.M{"$set": models.NewPartnerDocument(partner)})
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
		var doc models.PartnerDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		partners = append(partners, doc.ToEntity())
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return partners, nil
}
