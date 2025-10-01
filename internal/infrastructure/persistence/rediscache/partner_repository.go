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

type partnerRepository struct {
	repo   repositories.PartnerRepository
	client *goredis.Client
	ttl    time.Duration
}

func NewPartnerRepository(client *goredis.Client, ttl time.Duration, repo repositories.PartnerRepository) repositories.PartnerRepository {
	if client == nil || repo == nil {
		return repo
	}

	return &partnerRepository{
		repo:   repo,
		client: client,
		ttl:    mergeTTL(ttl, time.Minute),
	}
}

func (r *partnerRepository) GetPartnerByID(ctx context.Context, id primitive.ObjectID) (*entities.Partner, error) {
	if r == nil {
		return nil, nil
	}

	key := buildIDKey("partners", id.Hex())
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached entities.Partner
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=partners operation=get id=%s source=redis", id.Hex())
			return &cached, nil
		} else {
			log.Printf("cache: stale resource=partners operation=get id=%s error=%v", id.Hex(), unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	partner, err := r.repo.GetPartnerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if partner != nil {
		log.Printf("cache: miss resource=partners operation=get id=%s source=mongo", id.Hex())
		_ = r.savePartner(ctx, key, partner)
	} else {
		log.Printf("cache: miss resource=partners operation=get id=%s source=mongo result=empty", id.Hex())
	}

	return partner, nil
}

func (r *partnerRepository) CreatePartner(ctx context.Context, partner *entities.Partner) error {
	if err := r.repo.CreatePartner(ctx, partner); err != nil {
		return err
	}

	if partner != nil {
		_ = r.savePartner(ctx, buildIDKey("partners", partner.ID.Hex()), partner)
	}

	_ = invalidateResourceLists(ctx, r.client, "partners")

	return nil
}

func (r *partnerRepository) UpdatePartner(ctx context.Context, partner *entities.Partner) error {
	if err := r.repo.UpdatePartner(ctx, partner); err != nil {
		return err
	}

	if partner != nil && !partner.ID.IsZero() {
		_ = r.savePartner(ctx, buildIDKey("partners", partner.ID.Hex()), partner)
	}

	_ = invalidateResourceLists(ctx, r.client, "partners")

	return nil
}

func (r *partnerRepository) DeletePartner(ctx context.Context, id primitive.ObjectID) error {
	if err := r.repo.DeletePartner(ctx, id); err != nil {
		return err
	}

	_ = r.client.Del(ctx, buildIDKey("partners", id.Hex())).Err()
	_ = invalidateResourceLists(ctx, r.client, "partners")

	return nil
}

func (r *partnerRepository) ListPartners(ctx context.Context, filter map[string]interface{}) ([]*entities.Partner, error) {
	key := buildListKey("partners", filter)
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached []*entities.Partner
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=partners operation=list key=%s source=redis count=%d", key, len(cached))
			return cached, nil
		} else {
			log.Printf("cache: stale resource=partners operation=list key=%s error=%v", key, unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	partners, err := r.repo.ListPartners(ctx, filter)
	if err != nil {
		return nil, err
	}

	if partners != nil {
		if payload, marshalErr := json.Marshal(partners); marshalErr == nil {
			_ = r.client.Set(ctx, key, payload, r.ttl).Err()
		}
	}

	log.Printf("cache: miss resource=partners operation=list key=%s source=mongo count=%d", key, len(partners))

	return partners, nil
}

func (r *partnerRepository) savePartner(ctx context.Context, key string, partner *entities.Partner) error {
	if partner == nil {
		return nil
	}

	payload, err := json.Marshal(partner)
	if err != nil {
		return err
	}

	if err := r.client.Set(ctx, key, payload, r.ttl).Err(); err != nil {
		return err
	}

	return nil
}
