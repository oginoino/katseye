package config

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"katseye/internal/domain/services"
	webmiddleware "katseye/internal/infrastructure/web/middleware"
)

type MiddlewareSet struct {
	JWT gin.HandlerFunc
}

func buildMiddlewares(cfg AuthConfig, tokenService *services.TokenService) (MiddlewareSet, error) {
	set := MiddlewareSet{}

	secret := strings.TrimSpace(cfg.JWTSecret)
	if secret == "" {
		return set, fmt.Errorf("jwt secret is not configured")
	}

	options := []webmiddleware.JWTOption{webmiddleware.WithPublicPaths("/auth/login")}
	if tokenService != nil {
		options = append(options, webmiddleware.WithTokenRevocationChecker(tokenService))
	}

	middleware, err := webmiddleware.NewJWTAuthMiddleware(secret, options...)
	if err != nil {
		return set, fmt.Errorf("creating jwt middleware: %w", err)
	}

	set.JWT = middleware
	return set, nil
}

func (m MiddlewareSet) toRouterMiddlewares() []gin.HandlerFunc {
	var middlewares []gin.HandlerFunc
	if m.JWT != nil {
		middlewares = append(middlewares, m.JWT)
	}
	return middlewares
}
