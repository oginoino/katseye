package config

import (
	interfaces "katseye/internal/application/interfaces/repositories"
	mongorepositories "katseye/internal/infrastructure/persistence/mongodb/repositories"
	rediscache "katseye/internal/infrastructure/persistence/rediscache"
)

type RepositorySet struct {
	Product interfaces.ProductRepository
	Partner interfaces.PartnerRepository
	Address interfaces.AddressRepository
	User    interfaces.UserRepository
}

func buildRepositories(resources *MongoResources, cache *RedisResources) RepositorySet {
	if resources == nil {
		return RepositorySet{}
	}

	var productRepo interfaces.ProductRepository = mongorepositories.NewProductRepositoryMongo(resources.Collections.Products)
	var partnerRepo interfaces.PartnerRepository = mongorepositories.NewPartnerRepositoryMongo(resources.Collections.Partners)
	var addressRepo interfaces.AddressRepository = mongorepositories.NewAddressRepositoryMongo(resources.Collections.Addresses)
	var userRepo interfaces.UserRepository = mongorepositories.NewUserRepositoryMongo(resources.Collections.Users)

	if cache != nil && cache.Client != nil {
		productRepo = rediscache.NewProductRepository(cache.Client, cache.TTL, productRepo)
		partnerRepo = rediscache.NewPartnerRepository(cache.Client, cache.TTL, partnerRepo)
		addressRepo = rediscache.NewAddressRepository(cache.Client, cache.TTL, addressRepo)
		userRepo = rediscache.NewUserRepository(cache.Client, cache.TTL, userRepo)
	}

	return RepositorySet{
		Product: productRepo,
		Partner: partnerRepo,
		Address: addressRepo,
		User:    userRepo,
	}
}
