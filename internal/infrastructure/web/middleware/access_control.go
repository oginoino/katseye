package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"katseye/internal/domain/entities"
	"katseye/internal/infrastructure/web/response"
)

// RequireProfileTypes ensures the authenticated user has one of the allowed profile types.
// When no profile types are provided the middleware does not enforce any restriction.
func RequireProfileTypes(allowed ...entities.UserProfileType) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, profile := range allowed {
		normalized := strings.TrimSpace(strings.ToLower(profile.String()))
		if normalized == "" {
			continue
		}
		allowedSet[normalized] = struct{}{}
	}

	// When nothing is provided, allow any profile to continue.
	if len(allowedSet) == 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		rawClaims, exists := c.Get(contextKeyClaims)
		if !exists {
			response.NewForbiddenResponse(c, "Access denied", "profile information not available")
			c.Abort()
			return
		}

		claims, ok := rawClaims.(jwt.MapClaims)
		if !ok {
			response.NewForbiddenResponse(c, "Access denied", "invalid profile claims")
			c.Abort()
			return
		}

		rawProfileType, hasProfile := claims["profile_type"]
		if !hasProfile {
			response.NewForbiddenResponse(c, "Access denied", "profile type not present in token")
			c.Abort()
			return
		}

		profileType := strings.TrimSpace(strings.ToLower(fmt.Sprint(rawProfileType)))
		if profileType == "" {
			response.NewForbiddenResponse(c, "Access denied", "profile type not present in token")
			c.Abort()
			return
		}

		if _, allowed := allowedSet[profileType]; !allowed {
			response.NewForbiddenResponse(c, "Access denied", "insufficient profile permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
