package cache

import (
	"context"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type ConfigurationRepo struct{}

func New() *ConfigurationRepo {
	return &ConfigurationRepo{}
}

func (c *ConfigurationRepo) Get(ctx context.Context, deviceID string) (entity.ConfigurationResponse, error) {
	return entity.ConfigurationResponse{
		Configuration: entity.Configuration{
			HeartbeatPeriod: 10 * time.Second,
		},
	}, nil
}
