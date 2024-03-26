package resty

import (
	"context"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
)

// A simple cache interface
type Cache[T any] interface {
	// Get the item with the given key
	Get(context.Context, string) (T, error)
	// Set the value against the given key
	Set(context.Context, string, T) error
	// Expire all items from the cache
	Clear(context.Context) error
}

// generic wrapper around cache.Cache
type cacheImpl[T any] struct {
	cache *cache.Cache[T]
}

// Create a new Cache with the given ttl expiration for items
func NewCache[T any](ttl time.Duration) Cache[T] {
	gocacheClient := gocache.New(ttl, 10*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)
	return &cacheImpl[T]{cache: cache.New[T](gocacheStore)}
}

func (c *cacheImpl[T]) Get(ctx context.Context, key string) (T, error) {
	return c.cache.Get(ctx, key)
}

func (c *cacheImpl[T]) Set(ctx context.Context, key string, value T) error {
	return c.cache.Set(ctx, key, value)
}

func (c *cacheImpl[T]) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx)
}

// nilCache is a Cache which does nothing
type nilCache[T any] struct{}

func (n *nilCache[T]) Get(context.Context, string) (T, error) {
	return *new(T), nil
}

func (n *nilCache[T]) Set(context.Context, string, T) error {
	return nil
}

func (n *nilCache[T]) Clear(context.Context) error {
	return nil
}
