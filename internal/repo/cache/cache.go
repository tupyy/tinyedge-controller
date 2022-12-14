package cache

import (
	"context"

	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type MemCacheRepo struct {
	cache map[string]entity.ConfigurationResponse
}

func NewCacheRepo() *MemCacheRepo {
	return &MemCacheRepo{
		cache: make(map[string]entity.ConfigurationResponse),
	}
}

func (c *MemCacheRepo) Put(ctx context.Context, id string, confResponse entity.ConfigurationResponse) error {
	c.cache[id] = confResponse
	return nil
}

func (c *MemCacheRepo) Get(ctx context.Context, id string) (entity.ConfigurationResponse, error) {
	conf, found := c.cache[id]
	if !found {
		return entity.ConfigurationResponse{}, common.ErrResourceNotFound
	}
	return conf, nil
}

func (c *MemCacheRepo) Delete(ctx context.Context, id string) error {
	delete(c.cache, id)
	return nil
}
