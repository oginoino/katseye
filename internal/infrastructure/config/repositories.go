package config

import (
	"katseye/internal/domain/repositories"
	"katseye/internal/domain/security"
	mongorepositories "katseye/internal/infrastructure/persistence/mongodb/repositories"
	rediscache "katseye/internal/infrastructure/persistence/rediscache"
)

type RepositorySet struct {
	Product repositories.ProductRepository
	Partner repositories.PartnerRepository
	Address repositories.AddressRepository
	User    repositories.UserRepository
	Token   security.TokenStore
}

func buildRepositories(resources *MongoResources, cache *RedisResources) RepositorySet {
	if resources == nil {
		return RepositorySet{}
	}

	var productRepo repositories.ProductRepository = mongorepositories.NewProductRepositoryMongo(resources.Collections.Products)
	var partnerRepo repositories.PartnerRepository = mongorepositories.NewPartnerRepositoryMongo(resources.Collections.Partners)
	var addressRepo repositories.AddressRepository = mongorepositories.NewAddressRepositoryMongo(resources.Collections.Addresses)
	var userRepo repositories.UserRepository = mongorepositories.NewUserRepositoryMongo(resources.Collections.Users)
	var tokenStore security.TokenStore

	if cache != nil && cache.Client != nil {
		productRepo = rediscache.NewProductRepository(cache.Client, cache.TTL, productRepo)
		partnerRepo = rediscache.NewPartnerRepository(cache.Client, cache.TTL, partnerRepo)
		addressRepo = rediscache.NewAddressRepository(cache.Client, cache.TTL, addressRepo)
		userRepo = rediscache.NewUserRepository(cache.Client, cache.TTL, userRepo)
		tokenStore = rediscache.NewTokenStore(cache.Client)
	}

	return RepositorySet{
		Product: productRepo,
		Partner: partnerRepo,
		Address: addressRepo,
		User:    userRepo,
		Token:   tokenStore,
	}
}
