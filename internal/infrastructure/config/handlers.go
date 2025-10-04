package config

import (
	handlers "katseye/internal/infrastructure/web/handlers"
	webrouter "katseye/internal/infrastructure/web/router"
)

type HandlerSet struct {
	Product  *handlers.ProductHandler
	Partner  *handlers.PartnerHandler
	Address  *handlers.AddressHandler
	Consumer *handlers.ConsumerHandler
	Auth     *handlers.AuthHandler
}

func buildHandlers(services ServiceSet, authCfg AuthConfig) HandlerSet {
	handlerSet := HandlerSet{}

	if services.Product != nil {
		handlerSet.Product = handlers.NewProductHandler(services.Product, services.ProductTemplates)
	}

	if services.Partner != nil {
		handlerSet.Partner = handlers.NewPartnerHandler(services.Partner)
	}

	if services.Address != nil {
		handlerSet.Address = handlers.NewAddressHandler(services.Address)
	}

	if services.Consumer != nil {
		handlerSet.Consumer = handlers.NewConsumerHandler(services.Consumer)
	}

	if services.Auth != nil {
		handlerSet.Auth = handlers.NewAuthHandler(services.Auth, services.Token, services.Partner, services.Consumer, authCfg.JWTSecret)
	}

	return handlerSet
}

func (h HandlerSet) toRouterHandlers() webrouter.Handlers {
	return webrouter.Handlers{
		Product:  h.Product,
		Partner:  h.Partner,
		Address:  h.Address,
		Consumer: h.Consumer,
		Auth:     h.Auth,
	}
}
