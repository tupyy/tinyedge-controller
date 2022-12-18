package entity

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
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

func (c CertificateGroup) GetSerialNumber() string {
	if c.Certificate == nil {
		return ""
	}
	return fmt.Sprintf("%X", c.Certificate.SerialNumber)
}
