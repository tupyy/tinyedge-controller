package auth

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/tupyy/tinyedge-controller/internal/services/common"
)

type Service struct {
	certManager  common.CertificateReader
	deviceReader common.DeviceReader
}

func New(certManager common.CertificateReader, deviceReader common.DeviceReader) *Service {
	return &Service{certManager, deviceReader}
}

// Auth is a function that perfomers authentication.
// It verify that the certificate presented by in client is not revoked and it has the same serial number
// as the one signed by the operator during the registration phase.
// Auth is not applied on Enrol and Register methods.
func (s *Service) Auth(ctx context.Context, method string, deviceID string, certificate []byte) (context.Context, error) {
	if method == "Enrol" || method == "Register" {
		return ctx, nil
	}

	// get the device
	device, err := s.deviceReader.Get(ctx, deviceID)
	if err != nil {
		return ctx, fmt.Errorf("device %q not found", deviceID)
	}

	// get the real certificate
	realCertificate, err := s.certManager.GetCertificate(ctx, device.CertificateSerialNumber)
	if err != nil {
		return ctx, fmt.Errorf("unable to get device %q certificate with sn %q: %w", deviceID, device.CertificateSerialNumber, err)
	}

	presentedCertificate, err := s.decodeCertificate(certificate)
	if err != nil {
		return ctx, fmt.Errorf("unable to parse certificate: %w", err)
	}

	if realCertificate.Certificate.SerialNumber != presentedCertificate.SerialNumber {
		return ctx, fmt.Errorf("certificates don't match")
	}

	if realCertificate.IsRevoked {
		return ctx, fmt.Errorf("unable to authenticate device %q. The presented certificate is revoked.", deviceID)
	}

	return ctx, nil
}

func (s *Service) decodeCertificate(cert []byte) (*x509.Certificate, error) {
	decodedCertificate, _ := pem.Decode(cert)
	if decodedCertificate == nil {
		return nil, fmt.Errorf("unable to decode certificate to PEM")
	}

	certificate, err := x509.ParseCertificate(decodedCertificate.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse certficate: %w", err)
	}

	return certificate, nil
}
