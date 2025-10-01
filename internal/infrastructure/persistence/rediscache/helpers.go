package rediscache

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func mergeTTL(requested, fallback time.Duration) time.Duration {
	if requested <= 0 {
		return fallback
	}
	return requested
}

func deleteByPattern(ctx context.Context, client *goredis.Client, pattern string) error {
	if client == nil {
		return nil
	}

	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func buildListKey(resource string, filter map[string]interface{}) string {
	if len(filter) == 0 {
		return fmt.Sprintf("%s:list:all", resource)
	}

	keys := make([]string, 0, len(filter))
	for key := range filter {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, filter[key]))
	}

	return fmt.Sprintf("%s:list:%s", resource, strings.Join(parts, "&"))
}

func buildIDKey(resource, id string) string {
	return fmt.Sprintf("%s:id:%s", resource, id)
}

func invalidateResourceLists(ctx context.Context, client *goredis.Client, resource string) error {
	pattern := fmt.Sprintf("%s:list:*", resource)
	return deleteByPattern(ctx, client, pattern)
}
