package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	webrouter "katseye/internal/infrastructure/web/router"
)

type Application struct {
	Config       Config
	server       *webrouter.Server
	mongo        *MongoResources
	cache        *RedisResources
	repositories RepositorySet
	services     ServiceSet
	handlers     HandlerSet
	middlewares  MiddlewareSet
}

func Initialize() (*Application, error) {
	settings, err := Load()
	if err != nil {
		return nil, err
	}

	log.Printf(
		"config: environment=%s gin_mode=%s port=%s mongo_uri=%s mongo_database=%s jwt_secret_configured=%t",
		settings.Environment,
		settings.HTTP.GinMode,
		settings.HTTP.Port,
		settings.Mongo.URI,
		settings.Mongo.Database,
		strings.TrimSpace(settings.Auth.JWTSecret) != "",
	)

	mongoResources, err := newMongoResources(settings.Mongo)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongo: %w", err)
	}

	log.Printf("mongo: client initialized database=%s", settings.Mongo.Database)

	redisResources, err := newRedisResources(settings.Cache)
	if err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	if redisResources != nil {
		log.Printf("redis: cache enabled addr=%s db=%d ttl=%s", settings.Cache.Redis.Address, settings.Cache.Redis.DB, redisResources.TTL)
	} else {
		log.Printf("redis: cache disabled")
	}

	repositories := buildRepositories(mongoResources, redisResources)
	services := buildServices(repositories)
	handlers := buildHandlers(services, settings.Auth)
	middlewares, err := buildMiddlewares(settings.Auth)
	if err != nil {
		return nil, fmt.Errorf("configuring middlewares: %w", err)
	}
	server := buildHTTPServer(settings.HTTP, handlers.toRouterHandlers(), middlewares.toRouterMiddlewares())

	log.Printf("http: server configured port=%s mode=%s", settings.HTTP.Port, settings.HTTP.GinMode)

	return &Application{
		Config:       settings,
		server:       server,
		mongo:        mongoResources,
		cache:        redisResources,
		repositories: repositories,
		services:     services,
		handlers:     handlers,
		middlewares:  middlewares,
	}, nil
}

func (a *Application) Close(ctx context.Context) error {
	if a == nil {
		return nil
	}

	var firstErr error

	if a.mongo != nil {
		if err := a.mongo.Close(ctx); err != nil {
			firstErr = retainFirstError(firstErr, err)
		}
	}

	if a.cache != nil {
		if err := a.cache.Close(ctx); err != nil {
			firstErr = retainFirstError(firstErr, err)
		}
	}

	return firstErr
}

func (a *Application) RunHTTPServer() error {
	if a == nil || a.server == nil {
		return fmt.Errorf("http server not configured")
	}

	log.Printf("startup: listening on :%s (gin_mode=%s)", a.Config.HTTP.Port, a.Config.HTTP.GinMode)

	return a.server.Run()
}

func (a *Application) HTTPServer() *webrouter.Server {
	if a == nil {
		return nil
	}

	return a.server
}

func (a *Application) Repositories() RepositorySet {
	return a.repositories
}

func (a *Application) Services() ServiceSet {
	return a.services
}

func (a *Application) Handlers() HandlerSet {
	return a.handlers
}

func (a *Application) Middlewares() MiddlewareSet {
	return a.middlewares
}

func (a *Application) Cache() *RedisResources {
	if a == nil {
		return nil
	}

	return a.cache
}

func retainFirstError(current error, candidate error) error {
	if current != nil {
		return current
	}
	return candidate
}
