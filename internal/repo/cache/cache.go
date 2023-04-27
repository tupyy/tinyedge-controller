package cache

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
)

type MemCacheRepo struct {
	cache map[string]entity.DeviceConfiguration
}

func NewCacheRepo() *MemCacheRepo {
	return &MemCacheRepo{
		cache: make(map[string]entity.DeviceConfiguration),
	}
}

func (c *MemCacheRepo) Put(ctx context.Context, id string, confResponse entity.DeviceConfiguration) error {
	c.cache[id] = confResponse
	return nil
}

func (c *MemCacheRepo) Get(ctx context.Context, id string) (entity.DeviceConfiguration, error) {
	conf, found := c.cache[id]
	if !found {
		return entity.DeviceConfiguration{}, errService.NewResourceNotFoundError("configuration", id)
	}
	return conf, nil
}

func (c *MemCacheRepo) Delete(ctx context.Context, id string) error {
	delete(c.cache, id)
	return nil
}
