package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	vvault "github.com/hashicorp/vault/api"
	"github.com/tupyy/tinyedge-controller/internal/clients/vault"
	errService "github.com/tupyy/tinyedge-controller/internal/services/errors"
	"go.uber.org/zap"
)

type CertficateRepo struct {
	vault                *vault.Vault
	cnSuffix             string
	certificateMountPath string
	pkiRoleID            string
}

func NewCertificateRepository(v *vault.Vault, certificateMountPath, cnSuffix, pkiRoleID string) *CertficateRepo {
	return &CertficateRepo{
		vault:                v,
		certificateMountPath: certificateMountPath,
		cnSuffix:             cnSuffix,
		pkiRoleID:            pkiRoleID,
	}
}

func (c *CertficateRepo) GetCACertificate(ctx context.Context) ([]byte, error) {
	pathToRead := fmt.Sprintf("%s/cert/ca", c.certificateMountPath)

	secret, err := c.vault.Client.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return []byte{}, err
	}

	data := secret.Data["certificate"]
	certificate := bytes.NewBufferString(data.(string)).Bytes()

	return certificate, nil
}

func (c *CertficateRepo) GenerateCertificate(ctx context.Context, cn string, ttl time.Duration) ([]byte, []byte, []byte, error) {
	pathToRead := fmt.Sprintf("%s/issue/%s", c.certificateMountPath, c.pkiRoleID)

	data := map[string]interface{}{
		"common_name":        cn,
		"issuer_ref":         c.pkiRoleID,
		"ttl":                ttl.Seconds(),
		"format":             "pem",
		"ip_sans":            "127.0.0.1",
		"private_key_format": "pkcs8",
	}

	secret, err := c.vault.Client.Logical().WriteWithContext(ctx, pathToRead, data)
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	certificate, _ := extract(secret, "certificate")
	privateKey, _ := extract(secret, "private_key")
	ca, _ := extract(secret, "issuing_ca")

	zap.S().Debugw("certificate generated", "cn", cn, "ttl", ttl)

	return certificate, privateKey, ca, nil
}

func (c *CertficateRepo) GetCertificate(ctx context.Context, sn string) ([]byte, bool, time.Time, error) {
	pathToRead := fmt.Sprintf("%s/cert/%s", c.certificateMountPath, sn)

	secret, err := c.vault.Client.Logical().ReadWithContext(ctx, pathToRead)
	if err != nil {
		return []byte{}, false, time.Time{}, err
	}

	if secret == nil {
		return []byte{}, false, time.Time{}, errService.NewResourceNotFoundError("device certificate", sn)
	}

	certificate, _ := extract(secret, "certificate")
	rev := secret.Data["revocation_time"].(json.Number)
	revTime, _ := rev.Int64()
	var revocationTime time.Time
	if revTime != 0 {
		revocationTime = time.Unix(revTime, 0)
	}

	return certificate, revTime != 0, revocationTime, nil
}

func (c *CertficateRepo) SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) ([]byte, error) {
	pathToRead := fmt.Sprintf("%s/sign/%s", c.certificateMountPath, c.pkiRoleID)

	data := map[string]interface{}{
		"csr":         string(csr),
		"common_name": cn,
		"ttl":         ttl.String(),
	}

	secret, err := c.vault.Client.Logical().WriteWithContext(ctx, pathToRead, data)
	if err != nil {
		return []byte{}, err
	}

	certificate, _ := extract(secret, "certificate")

	zap.S().Debugw("certificate request signed", "cn", cn, "ttl", ttl)

	return certificate, nil
}

func extract(secret *vvault.Secret, key string) ([]byte, error) {
	data, ok := secret.Data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key %q not found", key)
	}
	return bytes.NewBufferString(data.(string)).Bytes(), nil
}
