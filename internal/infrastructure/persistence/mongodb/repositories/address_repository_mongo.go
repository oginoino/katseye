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

type AddressRepositoryMongo struct {
	collection *mongo.Collection
}

func NewAddressRepositoryMongo(collection *mongo.Collection) repositories.AddressRepository {
	return &AddressRepositoryMongo{
		collection: collection,
	}
}

func (r *AddressRepositoryMongo) GetAddressByID(ctx context.Context, id primitive.ObjectID) (*entities.Address, error) {
	var doc models.AddressDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Address not found
		}
		return nil, err
	}
	return doc.ToEntity(), nil
}

func (r *AddressRepositoryMongo) CreateAddress(ctx context.Context, address *entities.Address) error {
	_, err := r.collection.InsertOne(ctx, models.NewAddressDocument(address))
	return err
}

func (r *AddressRepositoryMongo) UpdateAddress(ctx context.Context, address *entities.Address) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": address.ID}, bson.M{"$set": models.NewAddressDocument(address)})
	return err
}

func (r *AddressRepositoryMongo) DeleteAddress(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *AddressRepositoryMongo) ListAddresses(ctx context.Context, filter map[string]interface{}) ([]*entities.Address, error) {
	var addresses []*entities.Address
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc models.AddressDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		addresses = append(addresses, doc.ToEntity())
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
