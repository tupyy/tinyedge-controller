package vault

import (
	"context"
)

type Client interface {
	GetSecret(ctx context.Context, enginePath, name, key string) (string, error)
}

type Repository struct {
	client     Client
	cnSuffix   string
	enginePath string
}

func New(v Client, enginePath string) *Repository {
	return &Repository{
		client:     v,
		enginePath: enginePath,
	}
}

func (r *Repository) GetSecret(ctx context.Context, name, key string) (string, error) {
	return r.client.GetSecret(ctx, r.enginePath, name, key)
}
