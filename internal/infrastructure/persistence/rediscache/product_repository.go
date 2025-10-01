package rediscache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"katseye/internal/application/interfaces/repositories"
	"katseye/internal/domain/entities"
)

type productRepository struct {
	repo   repositories.ProductRepository
	client *goredis.Client
	ttl    time.Duration
}

func NewProductRepository(client *goredis.Client, ttl time.Duration, repo repositories.ProductRepository) repositories.ProductRepository {
	if client == nil || repo == nil {
		return repo
	}

	return &productRepository{
		repo:   repo,
		client: client,
		ttl:    mergeTTL(ttl, time.Minute),
	}
}

func (r *productRepository) GetProductByID(ctx context.Context, id primitive.ObjectID) (*entities.Product, error) {
	if r == nil {
		return nil, nil
	}

	key := buildIDKey("products", id.Hex())
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached entities.Product
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=products operation=get id=%s source=redis", id.Hex())
			return &cached, nil
		} else {
			log.Printf("cache: stale resource=products operation=get id=%s error=%v", id.Hex(), unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	product, err := r.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product != nil {
		log.Printf("cache: miss resource=products operation=get id=%s source=mongo", id.Hex())
		_ = r.saveProduct(ctx, key, product)
	} else {
		log.Printf("cache: miss resource=products operation=get id=%s source=mongo result=empty", id.Hex())
	}

	return product, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product *entities.Product) error {
	if err := r.repo.CreateProduct(ctx, product); err != nil {
		return err
	}

	if product != nil {
		_ = r.saveProduct(ctx, buildIDKey("products", product.ID.Hex()), product)
	}

	_ = invalidateResourceLists(ctx, r.client, "products")

	return nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *entities.Product) error {
	if err := r.repo.UpdateProduct(ctx, product); err != nil {
		return err
	}

	if product != nil && !product.ID.IsZero() {
		_ = r.saveProduct(ctx, buildIDKey("products", product.ID.Hex()), product)
	}

	_ = invalidateResourceLists(ctx, r.client, "products")

	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	if err := r.repo.DeleteProduct(ctx, id); err != nil {
		return err
	}

	_ = r.client.Del(ctx, buildIDKey("products", id.Hex())).Err()
	_ = invalidateResourceLists(ctx, r.client, "products")

	return nil
}

func (r *productRepository) ListProducts(ctx context.Context, filter map[string]interface{}) ([]*entities.Product, error) {
	key := buildListKey("products", filter)
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached []*entities.Product
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=products operation=list key=%s source=redis count=%d", key, len(cached))
			return cached, nil
		} else {
			log.Printf("cache: stale resource=products operation=list key=%s error=%v", key, unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	products, err := r.repo.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	if products != nil {
		if payload, marshalErr := json.Marshal(products); marshalErr == nil {
			_ = r.client.Set(ctx, key, payload, r.ttl).Err()
		}
	}

	log.Printf("cache: miss resource=products operation=list key=%s source=mongo count=%d", key, len(products))

	return products, nil
}

func (r *productRepository) saveProduct(ctx context.Context, key string, product *entities.Product) error {
	if product == nil {
		return nil
	}

	payload, err := json.Marshal(product)
	if err != nil {
		return err
	}

	if err := r.client.Set(ctx, key, payload, r.ttl).Err(); err != nil {
		return err
	}

	return nil
}
