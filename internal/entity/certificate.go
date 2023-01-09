package entity

import (
	"crypto/x509"
	"fmt"
	"time"
)

type CertificateGroup struct {
	Certificate     *x509.Certificate
	PrivateKey      any
	CACertificate   *x509.Certificate
	CertificatePEM  []byte
	PrivateKeyPEM   []byte
	CACertficatePEM []byte
	RevocationTime  time.Time
	IsRevoked       bool
}

func (c CertificateGroup) GetSerialNumber() string {
	if c.Certificate == nil {
		return ""
	}
	return fmt.Sprintf("%X", c.Certificate.SerialNumber)
}
