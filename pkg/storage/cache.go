package storage

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	MGet(ctx context.Context, keys []string) (map[string]string, error)
	MSet(ctx context.Context, values map[string]string, expiration time.Duration) error
}

type CacheMock struct {
}

func (c *CacheMock) Get(ctx context.Context, key string) (string, error) {
	return key, nil
}

func (c *CacheMock) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return nil
}

func (c *CacheMock) MGet(ctx context.Context, keys []string) (map[string]string, error) {
	return nil, nil
}

func (c *CacheMock) MSet(ctx context.Context, values map[string]string, expiration time.Duration) error {
	return nil
}
