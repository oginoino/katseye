package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"katseye/internal/infrastructure/web/response"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	contextKeyToken     = "jwt_token"
	contextKeyClaims    = "jwt_claims"
	contextKeyRawToken  = "jwt_raw_token"
)

type jwtAuthConfig struct {
	publicPaths       map[string]struct{}
	revocationChecker TokenRevocationChecker
}

// JWTOption allows customizing the middleware behaviour.
type JWTOption func(*jwtAuthConfig)

// WithPublicPaths configures paths that bypass token validation.
func WithPublicPaths(paths ...string) JWTOption {
	return func(cfg *jwtAuthConfig) {
		if cfg.publicPaths == nil {
			cfg.publicPaths = make(map[string]struct{})
		}
		for _, p := range paths {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			cfg.publicPaths[p] = struct{}{}
		}
	}
}

// TokenRevocationChecker reports whether a token has been revoked.
type TokenRevocationChecker interface {
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}

// WithTokenRevocationChecker sets the component responsible for checking token revocation.
func WithTokenRevocationChecker(checker TokenRevocationChecker) JWTOption {
	return func(cfg *jwtAuthConfig) {
		cfg.revocationChecker = checker
	}
}

// NewJWTAuthMiddleware creates a Gin middleware that validates JWT bearer tokens using the provided secret.
func NewJWTAuthMiddleware(secret string, opts ...JWTOption) (gin.HandlerFunc, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, fmt.Errorf("jwt secret must not be empty")
	}

	config := &jwtAuthConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}

	signingKey := []byte(secret)

	return func(c *gin.Context) {
		if shouldSkipAuth(c, config) {
			c.Next()
			return
		}

		tokenString, err := extractBearerToken(c.GetHeader(authorizationHeader))
		if err != nil {
			response.NewUnauthorizedResponse(c, "Missing authorization token", err.Error())
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return signingKey, nil
		})
		if err != nil {
			response.NewUnauthorizedResponse(c, "Invalid token", err.Error())
			c.Abort()
			return
		}

		if !token.Valid {
			response.NewUnauthorizedResponse(c, "Invalid token", "token validation failed")
			c.Abort()
			return
		}

		if config.revocationChecker != nil {
			revoked, revocationErr := config.revocationChecker.IsTokenRevoked(c.Request.Context(), tokenString)
			if revocationErr != nil {
				response.NewInternalServerErrorResponse(c, "Token validation error", revocationErr.Error())
				c.Abort()
				return
			}
			if revoked {
				response.NewUnauthorizedResponse(c, "Invalid token", "token revoked")
				c.Abort()
				return
			}
		}

		c.Set(contextKeyToken, token)
		c.Set(contextKeyRawToken, tokenString)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set(contextKeyClaims, claims)
		}

		c.Next()
	}, nil
}

func extractBearerToken(header string) (string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", fmt.Errorf("authorization header required")
	}
	if !strings.HasPrefix(header, bearerPrefix) {
		return "", fmt.Errorf("authorization header must be in format 'Bearer <token>'")
	}
	token := strings.TrimSpace(header[len(bearerPrefix):])
	if token == "" {
		return "", fmt.Errorf("authorization token missing")
	}
	return token, nil
}

func shouldSkipAuth(c *gin.Context, cfg *jwtAuthConfig) bool {
	if c.Request.Method == http.MethodOptions {
		return true
	}
	if cfg == nil || len(cfg.publicPaths) == 0 {
		return false
	}

	if path := strings.TrimSpace(c.FullPath()); path != "" {
		if _, ok := cfg.publicPaths[path]; ok {
			return true
		}
	}

	if raw := strings.TrimSpace(c.Request.URL.Path); raw != "" {
		if _, ok := cfg.publicPaths[raw]; ok {
			return true
		}
	}

	return false
}
