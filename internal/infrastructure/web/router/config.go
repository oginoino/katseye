package router

import "github.com/gin-gonic/gin"

type Config struct {
	Port        string
	Mode        string
	Handlers    Handlers
	Middlewares []gin.HandlerFunc
}

func ConfigureRoutes(r gin.IRouter, h Handlers) {
	if r == nil {
		return
	}

	registerAuthRoutes(r, h.Auth)
	registerProductRoutes(r, h.Product)
	registerPartnerRoutes(r, h.Partner)
	registerAddressRoutes(r, h.Address)
	registerConsumerRoutes(r, h.Consumer)
}
