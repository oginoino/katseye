package config

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	webmiddleware "katseye/internal/infrastructure/web/middleware"
)

type MiddlewareSet struct {
	JWT gin.HandlerFunc
}

func buildMiddlewares(cfg AuthConfig) (MiddlewareSet, error) {
	set := MiddlewareSet{}

	secret := strings.TrimSpace(cfg.JWTSecret)
	if secret == "" {
		return set, fmt.Errorf("jwt secret is not configured")
	}

	middleware, err := webmiddleware.NewJWTAuthMiddleware(secret, webmiddleware.WithPublicPaths("/auth/login"))
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
