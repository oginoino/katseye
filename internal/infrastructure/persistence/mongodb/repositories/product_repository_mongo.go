package mongodb

import (
	"context"
	"katseye/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepositoryMongo struct {
	collection *mongo.Collection
}

func NewProductRepositoryMongo(collection *mongo.Collection) *productRepositoryMongo {
	return &productRepositoryMongo{
		collection: collection,
	}
}

func (r *productRepositoryMongo) GetProductByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	var product entities.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Product not found
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryMongo) CreateProduct(ctx context.Context, product *entities.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *productRepositoryMongo) UpdateProduct(ctx context.Context, product *entities.Product) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$set": product})
	return err
}

func (r *productRepositoryMongo) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *productRepositoryMongo) ListProducts(ctx context.Context, filter map[string]interface{}) ([]*entities.Product, error) {
	var products []*entities.Product
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product entities.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryMongo) Save(ctx context.Context, product *entities.Product) error {
	if product.ID.IsZero() {
		return r.CreateProduct(ctx, product)
	}
	return r.UpdateProduct(ctx, product)
}

func (r *productRepositoryMongo) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	return r.GetProductByID(ctx, id)
}

func (r *productRepositoryMongo) Update(ctx context.Context, product *entities.Product) error {
	return r.UpdateProduct(ctx, product)
}

func (r *productRepositoryMongo) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.DeleteProduct(ctx, id)
}

func (r *productRepositoryMongo) List(ctx context.Context, filter map[string]interface{}) ([]*entities.Product, error) {
	return r.ListProducts(ctx, filter)
}
