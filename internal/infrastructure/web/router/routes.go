package router

import (
	"github.com/gin-gonic/gin"
	"katseye/internal/infrastructure/web/handlers"
)

func registerAuthRoutes(r gin.IRouter, handler *handlers.AuthHandler) {
	if handler == nil {
		return
	}

	auth := r.Group("/auth")
	auth.POST("/login", handler.Login)
	auth.POST("/logout", handler.Logout)
	serviceAccounts := auth.Group("/service-accounts")
	serviceAccounts.POST("", handler.CreateUser)
	serviceAccounts.DELETE("/:id", handler.DeleteUser)
}

func registerProductRoutes(r gin.IRouter, handler *handlers.ProductHandler) {
	if handler == nil {
		return
	}

	products := r.Group("/products")
	products.GET("", handler.ListProducts)
	products.POST("", handler.CreateProduct)
	products.GET("/:id", handler.GetProduct)
	products.PUT("/:id", handler.UpdateProduct)
	products.DELETE("/:id", handler.DeleteProduct)
}

func registerPartnerRoutes(r gin.IRouter, handler *handlers.PartnerHandler) {
	if handler == nil {
		return
	}

	partners := r.Group("/partners")
	partners.GET("", handler.ListPartners)
	partners.POST("", handler.CreatePartner)
	partners.GET("/:id", handler.GetPartner)
	partners.PUT("/:id", handler.UpdatePartner)
	partners.DELETE("/:id", handler.DeletePartner)
}

func registerAddressRoutes(r gin.IRouter, handler *handlers.AddressHandler) {
	if handler == nil {
		return
	}

	addresses := r.Group("/addresses")
	addresses.GET("", handler.ListAddresses)
	addresses.POST("", handler.CreateAddress)
	addresses.GET("/:id", handler.GetAddress)
	addresses.PUT("/:id", handler.UpdateAddress)
	addresses.DELETE("/:id", handler.DeleteAddress)
}

func registerConsumerRoutes(r gin.IRouter, handler *handlers.ConsumerHandler) {
	if handler == nil {
		return
	}

	customers := r.Group("/customers")
	customers.GET("", handler.ListConsumers)
	customers.POST("", handler.CreateConsumer)
	customers.GET("/:id", handler.GetConsumer)
	customers.PUT("/:id", handler.UpdateConsumer)
	customers.DELETE("/:id", handler.DeleteConsumer)
	customers.POST("/:id/products/:product_id", handler.ContractProduct)
	customers.DELETE("/:id/products/:product_id", handler.RemoveProduct)
}
