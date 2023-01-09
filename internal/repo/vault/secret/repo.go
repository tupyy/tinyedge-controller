package vault

import (
	"bytes"
	"context"
	"crypto/sha512"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/clients/vault"
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type Repository struct {
	vault      *vault.Vault
	enginePath string
}

func New(v *vault.Vault, enginePath string) *Repository {
	return &Repository{
		vault:      v,
		enginePath: enginePath,
	}
}

func (r *Repository) GetSecret(ctx context.Context, path, key string) (entity.Secret, error) {
	secret, err := r.vault.Client.KVv2(r.enginePath).Get(ctx, path)
	if err != nil {
		return entity.Secret{}, fmt.Errorf("unable to read secret: %w", err)
	}

	data, ok := secret.Data[key]
	if !ok {
		return entity.Secret{}, fmt.Errorf("the secret retrieved from vault is missing %q field", key)
	}

	dataString, ok := data.(string)
	if !ok {
		return entity.Secret{}, fmt.Errorf("unexpected secret key type for %q field", key)
	}

	e := entity.Secret{
		Path:  path,
		Key:   key,
		Value: dataString,
		Hash:  r.compuateHash(path, key, dataString),
	}

	return e, nil
}

func (r *Repository) compuateHash(path, key, value string) string {
	hash := sha512.New()
	hash.Write(bytes.NewBufferString(fmt.Sprintf("%s%s%s", path, key, value)).Bytes())
	return fmt.Sprintf("%x", hash.Sum(nil))
}
