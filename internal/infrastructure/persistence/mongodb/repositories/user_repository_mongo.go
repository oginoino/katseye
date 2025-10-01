package mongodb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"katseye/internal/domain/entities"
)

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(collection *mongo.Collection) *UserRepositoryMongo {
	return &UserRepositoryMongo{collection: collection}
}

func (r *UserRepositoryMongo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	if r == nil || r.collection == nil {
		return nil, nil
	}

	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, nil
	}

	filter := bson.M{"email": email}
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "password_hash": 1, "email": 1, "active": 1})

	var user entities.User
	if err := r.collection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	user.Normalize()

	return &user, nil
}
