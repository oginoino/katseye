package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"

	defaultHTTPPort    = "8080"
	defaultMongoURI    = "mongodb://localhost:27017"
	defaultMongoDB     = "katseye"
	defaultRedisAddr   = "localhost:6379"
	defaultRedisDB     = 0
	defaultRedisTTL    = 5 * time.Minute
	defaultCORSOrigins = "*"
	defaultCORSMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	defaultCORSHeaders = "Authorization,Content-Type,Accept,Origin"

	redisEnabledEnvKey = "REDIS_ENABLED"
	redisAddrEnvKey    = "REDIS_ADDR"
	redisPasswordKey   = "REDIS_PASSWORD"
	redisDBEnvKey      = "REDIS_DB"
	redisTTLEnvKey     = "REDIS_CACHE_TTL"

	appEnvKey                  = "APP_ENV"
	ginModeEnvKey              = "GIN_MODE"
	jwtSecretEnvKey            = "JWT_SECRET"
	productionEnvFile          = ".env"
	developmentEnvFile         = ".env.example"
	corsAllowedOriginsEnvKey   = "CORS_ALLOWED_ORIGINS"
	corsAllowedMethodsEnvKey   = "CORS_ALLOWED_METHODS"
	corsAllowedHeadersEnvKey   = "CORS_ALLOWED_HEADERS"
	corsAllowCredentialsEnvKey = "CORS_ALLOW_CREDENTIALS"
)

type Config struct {
	Environment string
	HTTP        HTTPConfig
	Mongo       MongoConfig
	Auth        AuthConfig
	Cache       CacheConfig
}

type HTTPConfig struct {
	Port             string
	GinMode          string
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type MongoConfig struct {
	URI      string
	Database string
}

type AuthConfig struct {
	JWTSecret string
}

type CacheConfig struct {
	Enabled bool
	Redis   RedisConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
	TTL      time.Duration
}

var cachedConfig struct {
	once sync.Once
	cfg  Config
	err  error
}

func Load() (Config, error) {
	cachedConfig.once.Do(func() {
		envName := resolveEnvironment()
		if err := applyEnvironment(envName); err != nil {
			cachedConfig.err = err
			return
		}

		cachedConfig.cfg = Config{
			Environment: envName,
			HTTP: HTTPConfig{
				Port:             lookupEnv("PORT", defaultHTTPPort),
				GinMode:          lookupEnv(ginModeEnvKey, gin.ReleaseMode),
				AllowedOrigins:   parseCSV(lookupEnv(corsAllowedOriginsEnvKey, ""), defaultCORSOrigins),
				AllowedMethods:   parseCSV(lookupEnv(corsAllowedMethodsEnvKey, ""), defaultCORSMethods),
				AllowedHeaders:   parseCSV(lookupEnv(corsAllowedHeadersEnvKey, ""), defaultCORSHeaders),
				AllowCredentials: parseBool(lookupEnv(corsAllowCredentialsEnvKey, "")),
			},
			Mongo: MongoConfig{
				URI:      lookupEnv("MONGO_URI", defaultMongoURI),
				Database: lookupEnv("MONGO_DATABASE", defaultMongoDB),
			},
			Auth: AuthConfig{
				JWTSecret: lookupEnv(jwtSecretEnvKey, ""),
			},
			Cache: loadCacheConfig(),
		}
	})

	return cachedConfig.cfg, cachedConfig.err
}

func resolveEnvironment() string {
	if env := strings.TrimSpace(os.Getenv(appEnvKey)); env != "" {
		return env
	}

	if env := readEnvKey(developmentEnvFile, appEnvKey); env != "" {
		return env
	}

	if env := readEnvKey(productionEnvFile, appEnvKey); env != "" {
		return env
	}

	return EnvironmentDevelopment
}

func applyEnvironment(env string) error {
	if env == "" {
		env = EnvironmentProduction
	}

	if err := loadEnvFile(environmentFile(env)); err != nil {
		return err
	}

	if current := strings.TrimSpace(os.Getenv(appEnvKey)); current == "" {
		if err := os.Setenv(appEnvKey, env); err != nil {
			return fmt.Errorf("setting %s: %w", appEnvKey, err)
		}
	}

	return nil
}

func environmentFile(env string) string {
	switch strings.ToLower(env) {
	case EnvironmentProduction:
		return productionEnvFile
	default:
		return developmentEnvFile
	}
}

func loadEnvFile(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("env file %s is a directory", path)
	}

	if err := godotenv.Overload(path); err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}

	return nil
}

func readEnvKey(file, key string) string {
	values, err := godotenv.Read(file)
	if err != nil {
		return ""
	}

	value := strings.TrimSpace(values[key])
	return value
}

func lookupEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func loadCacheConfig() CacheConfig {
	cacheCfg := CacheConfig{}

	enabledValue := lookupEnv(redisEnabledEnvKey, "")
	cacheCfg.Enabled = parseBool(enabledValue)

	redisCfg := RedisConfig{
		Address:  lookupEnv(redisAddrEnvKey, defaultRedisAddr),
		Password: lookupEnv(redisPasswordKey, ""),
		DB:       parseInt(lookupEnv(redisDBEnvKey, ""), defaultRedisDB),
		TTL:      parseDuration(lookupEnv(redisTTLEnvKey, ""), defaultRedisTTL),
	}

	cacheCfg.Redis = redisCfg

	return cacheCfg
}

func parseBool(value string) bool {
	if value == "" {
		return false
	}
	value = strings.TrimSpace(value)
	switch strings.ToLower(value) {
	case "1", "true", "t", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func parseInt(value string, fallback int) int {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return fallback
}

func parseDuration(value string, fallback time.Duration) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return fallback
}

func parseCSV(value string, fallback string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return splitAndTrim(fallback)
	}
	return splitAndTrim(value)
}

func splitAndTrim(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{})
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, trimmed)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
