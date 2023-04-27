package vault

import (
	"bytes"
	"context"
	"crypto/sha512"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/clients/vault"
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type SecretRepository struct {
	vault      *vault.Vault
	enginePath string
}

func NewSecretRepository(v *vault.Vault, enginePath string) *SecretRepository {
	return &SecretRepository{
		vault:      v,
		enginePath: enginePath,
	}
}

func (r *SecretRepository) GetSecret(ctx context.Context, path, key string) (entity.Secret, error) {
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

func (r *SecretRepository) compuateHash(path, key, value string) string {
	hash := sha512.New()
	hash.Write(bytes.NewBufferString(fmt.Sprintf("%s%s%s", path, key, value)).Bytes())
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (r *SecretRepository) GetCredentialsFunc(ctx context.Context, authType entity.RepositoryAuthType, secretPath string) entity.CredentialsFunc {
	return func(ctx context.Context, path string) (interface{}, error) {
		secret, err := r.vault.Client.KVv2(r.enginePath).Get(ctx, path)
		if err != nil {
			return entity.Secret{}, fmt.Errorf("unable to read secret: %w", err)
		}

		switch authType {
		case entity.SSHRepositoryAuthType:
			// expect a key named privatekey
			data, ok := secret.Data["private_key"]
			password, pok := secret.Data["password"]
			if !ok {
				return nil, fmt.Errorf("SSH private key not found in secret %q", secretPath)
			}
			privateKey := bytes.NewBufferString(data.(string)).Bytes()
			pass := ""
			if pok {
				pass = password.(string)
			}
			return entity.SSHRepositoryAuth{
				Password:   pass,
				PrivateKey: privateKey,
			}, nil
		case entity.TokenRepositoryAuthType:
			token, ok := secret.Data["token"]
			if !ok {
				return nil, fmt.Errorf("Token not found in secret %q", secretPath)
			}
			return entity.TokenRepositoryAuth{Token: token.(string)}, nil
		case entity.BasicRepositoryAuthType:
			username, ok := secret.Data["username"]
			password, pok := secret.Data["password"]
			if !ok || !pok {
				return nil, fmt.Errorf("Either username or password not found in secret %q", secretPath)
			}
			return entity.BasicRepositoryAuth{Username: username.(string), Password: password.(string)}, nil
		}
		return nil, fmt.Errorf("unknown repository auth type")
	}
}
