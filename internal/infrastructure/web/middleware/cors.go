package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

func NewCORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	origins, allowAllOrigins := normalizeOrigins(cfg.AllowedOrigins)
	methods := normalizeMethods(cfg.AllowedMethods)
	headers, allowAllHeaders := normalizeHeaders(cfg.AllowedHeaders)

	allowedMethodsValue := strings.Join(methods, ", ")
	allowedHeadersValue := strings.Join(headers, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		writerHeaders := c.Writer.Header()
		writerHeaders.Add("Vary", "Origin")
		if allowAllHeaders {
			writerHeaders.Add("Vary", "Access-Control-Request-Headers")
		}

		originAllowed := isOriginAllowed(origin, origins, allowAllOrigins)
		if originAllowed {
			switch {
			case cfg.AllowCredentials && origin != "":
				writerHeaders.Set("Access-Control-Allow-Origin", origin)
				writerHeaders.Set("Access-Control-Allow-Credentials", "true")
			case allowAllOrigins:
				writerHeaders.Set("Access-Control-Allow-Origin", "*")
			case origin != "":
				writerHeaders.Set("Access-Control-Allow-Origin", origin)
			}

			if allowedMethodsValue != "" {
				writerHeaders.Set("Access-Control-Allow-Methods", allowedMethodsValue)
			}

			if allowAllHeaders {
				if requestHeaders := strings.TrimSpace(c.GetHeader("Access-Control-Request-Headers")); requestHeaders != "" {
					writerHeaders.Set("Access-Control-Allow-Headers", requestHeaders)
				} else {
					writerHeaders.Set("Access-Control-Allow-Headers", "*")
				}
			} else if allowedHeadersValue != "" {
				writerHeaders.Set("Access-Control-Allow-Headers", allowedHeadersValue)
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func normalizeOrigins(origins []string) ([]string, bool) {
	if len(origins) == 0 {
		return nil, true
	}

	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(origins))
	allowAll := false

	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			allowAll = true
			continue
		}

		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	if allowAll {
		return normalized, true
	}

	if len(normalized) == 0 {
		return nil, false
	}

	return normalized, false
}

func normalizeMethods(methods []string) []string {
	defaults := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}

	if len(methods) == 0 {
		return defaults
	}

	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(methods))

	for _, method := range methods {
		trimmed := strings.TrimSpace(method)
		if trimmed == "" {
			continue
		}
		upper := strings.ToUpper(trimmed)
		if _, exists := seen[upper]; exists {
			continue
		}
		seen[upper] = struct{}{}
		normalized = append(normalized, upper)
	}

	if len(normalized) == 0 {
		return defaults
	}

	if _, ok := seen[http.MethodOptions]; !ok {
		normalized = append(normalized, http.MethodOptions)
	}

	return normalized
}

func normalizeHeaders(headers []string) ([]string, bool) {
	defaults := []string{"Authorization", "Content-Type", "Accept", "Origin"}

	if len(headers) == 0 {
		return defaults, false
	}

	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(headers))
	wildcard := false

	for _, header := range headers {
		trimmed := strings.TrimSpace(header)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			wildcard = true
			continue
		}

		canonical := http.CanonicalHeaderKey(trimmed)
		key := strings.ToLower(canonical)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, canonical)
	}

	if wildcard {
		return normalized, true
	}

	if len(normalized) == 0 {
		return defaults, false
	}

	return normalized, false
}

func isOriginAllowed(origin string, allowed []string, allowAll bool) bool {
	if origin == "" {
		return true
	}
	if allowAll {
		return true
	}
	lowerOrigin := strings.ToLower(origin)
	for _, allowedOrigin := range allowed {
		if strings.ToLower(allowedOrigin) == lowerOrigin {
			return true
		}
	}
	return false
}
