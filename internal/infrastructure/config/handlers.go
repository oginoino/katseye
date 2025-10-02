package config

import (
	handlers "katseye/internal/infrastructure/web/handlers"
	webrouter "katseye/internal/infrastructure/web/router"
)

type HandlerSet struct {
	Product *handlers.ProductHandler
	Partner *handlers.PartnerHandler
	Address *handlers.AddressHandler
	Auth    *handlers.AuthHandler
}

func buildHandlers(services ServiceSet, authCfg AuthConfig) HandlerSet {
	handlerSet := HandlerSet{}

	if services.Product != nil {
		handlerSet.Product = handlers.NewProductHandler(services.Product)
	}

	if services.Partner != nil {
		handlerSet.Partner = handlers.NewPartnerHandler(services.Partner)
	}

	if services.Address != nil {
		handlerSet.Address = handlers.NewAddressHandler(services.Address)
	}

	if services.Auth != nil {
		handlerSet.Auth = handlers.NewAuthHandler(services.Auth, services.Token, authCfg.JWTSecret)
	}

	return handlerSet
}

func (h HandlerSet) toRouterHandlers() webrouter.Handlers {
	return webrouter.Handlers{
		Product: h.Product,
		Partner: h.Partner,
		Address: h.Address,
		Auth:    h.Auth,
	}
}
