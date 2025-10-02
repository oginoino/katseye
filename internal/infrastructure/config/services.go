package config

import "katseye/internal/domain/services"

type ServiceSet struct {
	Product *services.ProductService
	Partner *services.PartnerService
	Address *services.AddressService
	Auth    *services.AuthService
	Token   *services.TokenService
}

func buildServices(repos RepositorySet) ServiceSet {
	return ServiceSet{
		Product: services.NewProductService(repos.Product, repos.Partner),
		Partner: services.NewPartnerService(repos.Partner),
		Address: services.NewAddressService(repos.Address),
		Auth:    services.NewAuthService(repos.User),
		Token:   services.NewTokenService(repos.Token),
	}
}
