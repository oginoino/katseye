# Internal Package

This directory contains the core implementation of the Katseye API application, organized following Domain-Driven Design (DDD) principles.

## Directory Structure

```
internal/
├── application/         # Application layer - Use cases and DTOs
├── domain/              # Domain layer - Business logic and rules
│   ├── entities/        # Domain entities
│   ├── repositories/    # Repository interfaces
│   ├── security/        # Security-related domain components
│   ├── services/        # Domain services
│   └── value_objects/   # Value objects
├── infrastructure/      # Infrastructure layer - External dependencies
│   ├── cache/           # Caching implementations
│   ├── config/          # Application configuration
│   ├── persistence/     # Database implementations
│   └── web/             # Web-related components
└── shared/              # Shared utilities and helpers
```

## Layers

### Application Layer

Contains application-specific logic, use cases, and DTOs (Data Transfer Objects). This layer orchestrates the flow of data between the domain layer and external interfaces.

### Domain Layer

The core of the application containing business logic and rules, independent of external concerns.

#### Entities

Domain entities representing the core business objects:
- `address.go` - Address entity
- `consumer.go` - Consumer entity
- `partner.go` - Partner entity
- `product.go` - Product entity
- `user.go` - User entity

#### Repositories

Repository interfaces defining data access contracts:
- `address_repository.go` - Address repository interface
- `consumer_repository.go` - Consumer repository interface
- `partner_repository.go` - Partner repository interface
- `product_repository.go` - Product repository interface
- `user_repository.go` - User repository interface

#### Services

Domain services implementing business logic:
- `address_service.go` - Address-related business logic
- `auth_service.go` - Authentication service
- `consumer_service.go` - Consumer-related business logic
- `partner_service.go` - Partner-related business logic
- `product_service.go` - Product-related business logic
- `token_service.go` - Token management service

#### Value Objects

Immutable objects that represent domain concepts:
- `address_type.go` - Types of addresses
- `consumer_type.go` - Types of consumers
- `partner_type.go` - Types of partners
- `product_category.go` - Product categories
- `product_type.go` - Types of products
- `required_document.go` - Required document specifications

### Infrastructure Layer

Implements interfaces defined in the domain layer and provides concrete implementations for external dependencies.

#### Cache

Caching implementations:
- `redis/` - Redis cache implementation

#### Config

Application configuration:
- `application.go` - Main application configuration
- `config.go` - Configuration loading
- `handlers.go` - Handler configuration
- `http.go` - HTTP server configuration
- `middleware.go` - Middleware configuration
- `mongo.go` - MongoDB configuration
- `redis.go` - Redis configuration
- `repositories.go` - Repository configuration
- `services.go` - Service configuration

#### Persistence

Database implementations:
- `mongodb/` - MongoDB implementations
- `rediscache/` - Redis cache implementations

#### Web

Web-related components:
- `dto/` - Data Transfer Objects for web layer
- `handlers/` - HTTP request handlers
- `middleware/` - HTTP middleware
- `response/` - HTTP response utilities
- `router/` - HTTP routing

### Shared Layer

Contains utilities and helpers shared across different layers of the application.