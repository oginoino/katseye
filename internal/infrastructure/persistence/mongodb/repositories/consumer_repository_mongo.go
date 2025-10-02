package mongodb

import (
	"context"

	"katseye/internal/domain/entities"
	"katseye/internal/infrastructure/persistence/mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type consumerRepositoryMongo struct {
	collection *mongo.Collection
}

func NewConsumerRepositoryMongo(collection *mongo.Collection) *consumerRepositoryMongo {
	return &consumerRepositoryMongo{collection: collection}
}

func (r *consumerRepositoryMongo) GetConsumerByID(ctx context.Context, id primitive.ObjectID) (*entities.Consumer, error) {
	var doc models.ConsumerDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return doc.ToEntity(), nil
}

func (r *consumerRepositoryMongo) CreateConsumer(ctx context.Context, consumer *entities.Consumer) error {
	_, err := r.collection.InsertOne(ctx, models.NewConsumerDocument(consumer))
	return err
}

func (r *consumerRepositoryMongo) UpdateConsumer(ctx context.Context, consumer *entities.Consumer) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": consumer.ID},
		bson.M{"$set": models.NewConsumerDocument(consumer)},
	)
	return err
}

func (r *consumerRepositoryMongo) DeleteConsumer(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *consumerRepositoryMongo) ListConsumers(ctx context.Context, filter map[string]interface{}) ([]*entities.Consumer, error) {
	var consumers []*entities.Consumer

	var bsonFilter interface{} = bson.D{}
	if len(filter) > 0 {
		m := bson.M{}
		for k, v := range filter {
			m[k] = v
		}
		bsonFilter = m
	}

	cursor, err := r.collection.Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc models.ConsumerDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		consumers = append(consumers, doc.ToEntity())
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return consumers, nil
}
