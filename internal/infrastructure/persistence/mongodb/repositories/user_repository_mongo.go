package mongodb

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"
	"katseye/internal/infrastructure/persistence/mongodb/models"
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
	opts := options.FindOne().SetProjection(bson.M{
		"_id":           1,
		"password_hash": 1,
		"email":         1,
		"active":        1,
		"role":          1,
		"permissions":   1,
	})

	var doc models.UserDocument
	if err := r.collection.FindOne(ctx, filter, opts).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return doc.ToEntity(), nil
}

func (r *UserRepositoryMongo) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	if r == nil || r.collection == nil {
		return nil, nil
	}

	if id.IsZero() {
		return nil, nil
	}

	filter := bson.M{"_id": id}
	opts := options.FindOne().SetProjection(bson.M{
		"_id":           1,
		"password_hash": 1,
		"email":         1,
		"active":        1,
		"role":          1,
		"permissions":   1,
	})

	var doc models.UserDocument
	if err := r.collection.FindOne(ctx, filter, opts).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return doc.ToEntity(), nil
}

func (r *UserRepositoryMongo) CreateUser(ctx context.Context, user *entities.User) error {
	if r == nil || r.collection == nil {
		return errors.New("user repository not configured")
	}
	if user == nil {
		return errors.New("user payload must not be nil")
	}

	_, err := r.collection.InsertOne(ctx, models.NewUserDocument(user))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return repositories.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepositoryMongo) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	if r == nil || r.collection == nil {
		return errors.New("user repository not configured")
	}
	if id.IsZero() {
		return repositories.ErrUserNotFound
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return repositories.ErrUserNotFound
	}

	return nil
}
