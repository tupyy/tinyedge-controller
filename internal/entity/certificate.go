package entity

import (
	"crypto/rsa"
	"crypto/x509"
	"time"
)

type CertificateGroup struct {
	Certificate    *x509.Certificate
	PrivateKey     *rsa.PrivateKey
	CertificatePEM []byte
	PrivateKeyPEM  []byte
	RevocationTime time.Time
	IsRevoked      bool
}
