package auth

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/tupyy/tinyedge-controller/internal/services/certificate"
	"go.uber.org/zap"
)

type Service struct {
	certManager  *certificate.Service
	deviceReader DeviceReader
}

func New(certManager *certificate.Service, deviceReader DeviceReader) *Service {
	return &Service{certManager, deviceReader}
}

// Auth is a function that perfomers authentication.
// It verify that the certificate presented by in client is not revoked and it has the same serial number
// as the one signed by the operator during the registration phase.
// Auth is not applied on Enrol and Register methods.
func (s *Service) Auth(ctx context.Context, method string, deviceID string, peerCertificates []*x509.Certificate) (context.Context, error) {
	// put the device id into context to be propagated though layers
	newCtx := context.WithValue(ctx, "device_id", deviceID)

	presentedCertificate := peerCertificates[0]

	if method == "/EdgeService/Enrol" || method == "/EdgeService/Register" {
		// we allow *only* registration certificate to authenticate to Register and Enrol methods
		if !s.isRegistrationCerficate(presentedCertificate) {
			return newCtx, fmt.Errorf("unable to authenticated the device %q for Enrol/Register methods. The certificate is invalid", deviceID)
		}
		return newCtx, nil
	}

	// registration device is not allowed from here
	if s.isRegistrationCerficate(presentedCertificate) {
		return newCtx, fmt.Errorf("unable to authenticated the device %q. It is forbidden to access method %q with a registration certificate", deviceID, method)
	}

	// get the device
	device, err := s.deviceReader.GetDevice(ctx, deviceID)
	if err != nil {
		return newCtx, fmt.Errorf("device %q not found", deviceID)
	}

	// get the real certificate
	realCertificate, err := s.certManager.GetCertificate(ctx, device.CertificateSerialNumber)
	if err != nil {
		return newCtx, fmt.Errorf("unable to get device %q certificate with sn %q: %w", deviceID, device.CertificateSerialNumber, err)
	}

	if s.getSerialNumber(realCertificate.Certificate) != s.getSerialNumber(presentedCertificate) {
		zap.S().Errorw("unable to authenticate. certificates don't match",
			"device_id", deviceID,
			"method", method,
			"presented_certificate_sn", s.getSerialNumber(presentedCertificate),
			"presented_certificate_cn", presentedCertificate.Subject.CommonName,
			"device_certificate_sn", s.getSerialNumber(realCertificate.Certificate),
			"device_certificate_cn", realCertificate.Certificate.Subject.CommonName,
		)
		return newCtx, fmt.Errorf("certificates don't match")
	}

	if realCertificate.IsRevoked {
		zap.S().Errorw("unable to authenticate device. the certificate is revoked",
			"device_id", deviceID,
			"method", method,
			"certificate_sn", s.getSerialNumber(realCertificate.Certificate),
			"certificate_cn", realCertificate.Certificate.Subject.CommonName,
			"certificate_revocation_time", realCertificate.RevocationTime,
		)
		return newCtx, fmt.Errorf("unable to authenticate device %q. The presented certificate is revoked.", deviceID)
	}

	zap.S().Infow("device authenticated", "method", method, "device_id", deviceID, "certificate_sn", s.getSerialNumber(realCertificate.Certificate))

	return newCtx, nil
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

func (s *Service) isRegistrationCerficate(cert *x509.Certificate) bool {
	return strings.HasPrefix(cert.Subject.CommonName, "registration")
}

func (s *Service) getSerialNumber(cert *x509.Certificate) string {
	return fmt.Sprintf("%x", cert.SerialNumber)
}
