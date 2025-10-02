package rediscache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"
)

type addressRepository struct {
	repo   repositories.AddressRepository
	client *goredis.Client
	ttl    time.Duration
}

func NewAddressRepository(client *goredis.Client, ttl time.Duration, repo repositories.AddressRepository) repositories.AddressRepository {
	if client == nil || repo == nil {
		return repo
	}

	return &addressRepository{
		repo:   repo,
		client: client,
		ttl:    mergeTTL(ttl, time.Minute),
	}
}

func (r *addressRepository) GetAddressByID(ctx context.Context, id primitive.ObjectID) (*entities.Address, error) {
	if r == nil {
		return nil, nil
	}

	key := buildIDKey("addresses", id.Hex())
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached entities.Address
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=addresses operation=get id=%s source=redis", id.Hex())
			return &cached, nil
		} else {
			log.Printf("cache: stale resource=addresses operation=get id=%s error=%v", id.Hex(), unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	address, err := r.repo.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if address != nil {
		log.Printf("cache: miss resource=addresses operation=get id=%s source=mongo", id.Hex())
		_ = r.saveAddress(ctx, key, address)
	} else {
		log.Printf("cache: miss resource=addresses operation=get id=%s source=mongo result=empty", id.Hex())
	}

	return address, nil
}

func (r *addressRepository) CreateAddress(ctx context.Context, address *entities.Address) error {
	if err := r.repo.CreateAddress(ctx, address); err != nil {
		return err
	}

	if address != nil {
		_ = r.saveAddress(ctx, buildIDKey("addresses", address.ID.Hex()), address)
	}

	_ = invalidateResourceLists(ctx, r.client, "addresses")

	return nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, address *entities.Address) error {
	if err := r.repo.UpdateAddress(ctx, address); err != nil {
		return err
	}

	if address != nil && !address.ID.IsZero() {
		_ = r.saveAddress(ctx, buildIDKey("addresses", address.ID.Hex()), address)
	}

	_ = invalidateResourceLists(ctx, r.client, "addresses")

	return nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, id primitive.ObjectID) error {
	if err := r.repo.DeleteAddress(ctx, id); err != nil {
		return err
	}

	_ = r.client.Del(ctx, buildIDKey("addresses", id.Hex())).Err()
	_ = invalidateResourceLists(ctx, r.client, "addresses")

	return nil
}

func (r *addressRepository) ListAddresses(ctx context.Context, filter map[string]interface{}) ([]*entities.Address, error) {
	key := buildListKey("addresses", filter)
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached []*entities.Address
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=addresses operation=list key=%s source=redis count=%d", key, len(cached))
			return cached, nil
		} else {
			log.Printf("cache: stale resource=addresses operation=list key=%s error=%v", key, unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	addresses, err := r.repo.ListAddresses(ctx, filter)
	if err != nil {
		return nil, err
	}

	if addresses != nil {
		if payload, marshalErr := json.Marshal(addresses); marshalErr == nil {
			_ = r.client.Set(ctx, key, payload, r.ttl).Err()
		}
	}

	log.Printf("cache: miss resource=addresses operation=list key=%s source=mongo count=%d", key, len(addresses))

	return addresses, nil
}

func (r *addressRepository) saveAddress(ctx context.Context, key string, address *entities.Address) error {
	if address == nil {
		return nil
	}

	payload, err := json.Marshal(address)
	if err != nil {
		return err
	}

	if err := r.client.Set(ctx, key, payload, r.ttl).Err(); err != nil {
		return err
	}

	return nil
}
