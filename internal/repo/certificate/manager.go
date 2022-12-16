package certificate

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tupyy/tinyedge-controller/internal/clients/vault"
	"github.com/tupyy/tinyedge-controller/internal/entity"
)

type Client interface {
	GetCACertificate(ctx context.Context) ([]byte, error)
	GenerateCertificate(ctx context.Context, cn string, ttl time.Duration) ([]byte, []byte, error)
	SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) ([]byte, error)
}

type Manager struct {
	vaultClient          *vault.Vault
	cnSuffix             string
	certificateMountPath string
}

func New(v *vault.Vault, certificateMountPath, cnSuffix string) *Manager {
	return &Manager{
		vaultClient:          v,
		certificateMountPath: certificateMountPath,
		cnSuffix:             cnSuffix,
	}
}

func (m *Manager) GetCACertificate(ctx context.Context) (entity.CertificateGroup, error) {
	certificate, err := m.vaultClient.GetCACertificate(ctx)
	if err != nil {
		return entity.CertificateGroup{}, err
	}

	decodedCertficate, _ := pem.Decode(certificate)
	cert, err := x509.ParseCertificate(decodedCertficate.Bytes)
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to parse certificate: %w", err)
	}

	return entity.CertificateGroup{
		Certificate:    cert,
		CertificatePEM: decodedCertficate.Bytes,
	}, nil
}

func (m *Manager) GetCertificate(ctx context.Context, serialNumber string) (entity.CertificateGroup, error) {
	formatSerialNumber := func(sn string) string {
		var sb strings.Builder
		for i := 2; true; i += 2 {
			if i >= len(serialNumber) {
				fmt.Fprintf(&sb, "%s", sn[i-2:i])
				break
			} else {
				fmt.Fprintf(&sb, "%s:", sn[i-2:i])
			}
		}
		return sb.String()
	}

	cert, isRevoked, revTime, err := m.vaultClient.GetCertificate(ctx, formatSerialNumber(serialNumber))
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to read certificate %q", serialNumber)
	}

	certificate, err := m.decode(cert)
	if err != nil {
		return entity.CertificateGroup{}, err
	}

	certificate.IsRevoked = isRevoked
	certificate.RevocationTime = revTime

	return certificate, nil
}

// GetServerCertificate returns the certificate used in mTLS.
func (m *Manager) GetServerCertificate(ctx context.Context, ttl time.Duration) (entity.CertificateGroup, error) {
	cn := "operator.home.net"
	hostname, err := os.Hostname()
	if err == nil {
		cn = fmt.Sprintf("%s-%s", hostname, cn)
	}
	return m.generateCertificate(ctx, cn, ttl)
}

// GenerateRegistrationCertificate returns a certificate used by the agent to registered itself.
func (m *Manager) GenerateRegistrationCertificate(ctx context.Context, ttl time.Duration) (entity.CertificateGroup, error) {
	return m.generateCertificate(ctx, "register.home.net", ttl)
}

func (m *Manager) SignCSR(ctx context.Context, csr []byte, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
	cert, err := m.vaultClient.SignCSR(ctx, csr, cn, ttl)
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to sign csr: %w", err)
	}

	return m.decode(cert)
}

func (m *Manager) generateCertificate(ctx context.Context, cn string, ttl time.Duration) (entity.CertificateGroup, error) {
	certificate, privateKey, err := m.vaultClient.GenerateCertificate(ctx, cn, ttl)

	decodedCertficate, _ := pem.Decode(certificate)
	cert, err := x509.ParseCertificate(decodedCertficate.Bytes)
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to parse certficate: %w", err)
	}

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to decode private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to parse private key: %w", err)
	}

	return entity.CertificateGroup{
		Certificate:    cert,
		PrivateKey:     key,
		CertificatePEM: decodedCertficate.Bytes,
		PrivateKeyPEM:  block.Bytes,
	}, nil
}

func (m *Manager) decode(cert []byte) (entity.CertificateGroup, error) {
	decodedCertificate, _ := pem.Decode(cert)
	if decodedCertificate == nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to decode certificate to PEM")
	}

	certificate, err := x509.ParseCertificate(decodedCertificate.Bytes)
	if err != nil {
		return entity.CertificateGroup{}, fmt.Errorf("unable to parse certficate: %w", err)
	}

	pemBlock := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: decodedCertificate.Bytes,
	})

	return entity.CertificateGroup{
		Certificate:    certificate,
		CertificatePEM: pemBlock,
	}, nil
}
