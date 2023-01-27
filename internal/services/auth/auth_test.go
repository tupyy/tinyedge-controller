package auth_test

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/internal/services/auth"
)

const (
	registrationCertificate = `
-----BEGIN CERTIFICATE-----
MIIDczCCAlugAwIBAgIUKSa6IUtOPaKLGqktgwU6Lj0pL9wwDQYJKoZIhvcNAQEL
BQAwKjEoMCYGA1UEAxMfaG9tZS5uZXQgSW50ZXJtZWRpYXRlIEF1dGhvcml0eTAe
Fw0yMzAxMjcxNTEwMjNaFw0yMzAyMjYxNTEwNTNaMCAxHjAcBgNVBAMTFXJlZ2lz
dHJhdGlvbi5ob21lLm5ldDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AKNxiMlB7y3D3OY415QstwHnpRFArGsXoi9a5qGGWOqZxHfC3L9y4p5yPRl8r32n
weRM74nxzYHv2wVZuZyEDH+UFJejoxVSXd0QxZ8kokI3gWCWTqcIyNWV5DPPAbVH
hCQA3iCo0HtMFVsWmK6Cd40bzVjM1F4NX7ja4QNNr+50Ku5yE2jnrxc1cezCXB4A
uiYPTgNEnaBD5VH8vgUVxbWJuuo+A5RD6dNhkrFoL7WhMMDw/5kRO6C2dawn5AV6
6sfadbMd1Ecx8BHWuF0463pIiLotmMrJFZHNtL3hpVh+cZTGX9gBOPC1ef3jDZH/
EpRX6JJi2a0icd1dnkWSdGMCAwEAAaOBmjCBlzAOBgNVHQ8BAf8EBAMCA6gwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQWBBTeX9M7KEcRdtm7
uBRtxQyhBySZITAfBgNVHSMEGDAWgBT/cHRmICNafdJm7uaJ6nVADziO9jAmBgNV
HREEHzAdghVyZWdpc3RyYXRpb24uaG9tZS5uZXSHBH8AAAEwDQYJKoZIhvcNAQEL
BQADggEBAE727SZoDqINh6g3fp03syl5bHNTunzCsRaFUCavIv7IGV5VOS7wDsZN
hxXCT0mJ/rRUfEu05qLQFb0rpBAs1eNxltMflt+YYNY35iD8K3GPlE3Ktl6a5Dam
ZN6xxgZrOSJo7fPotYfmrTnp+jxd73751L1qgQzb8cuxX1tS47AJwROg48u5EwCZ
HHAEOSxG0wodyKYfXP+VmCm4sMBZtMUGCrZ3OHAVbc8KKx4kEoWF6v0LgA6Fbd98
dXJjC078QlcTjCSiUMWnWLivUsttf1hPGTAk/64fBDWX9GqXYWt5L+gxfUXMK8gN
mASdl8bXVXLr0pk7oRjo0cL2z5nzJ1s=
-----END CERTIFICATE-----
`

	clientCertificate = `
-----BEGIN CERTIFICATE-----
MIIDZzCCAk+gAwIBAgIUYJTeyIop/BahgNnfBYk8SC5t+i0wDQYJKoZIhvcNAQEL
BQAwKjEoMCYGA1UEAxMfaG9tZS5uZXQgSW50ZXJtZWRpYXRlIEF1dGhvcml0eTAe
Fw0yMzAxMjcxNTExMTdaFw0yMzAyMjYxNTExNDZaMBoxGDAWBgNVBAMTD3NlcnZl
ci5ob21lLm5ldDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMd5YirX
boiTKK1+PhBNXqrPgqAGxaWHyEtu1kOpsqsSoir5xXjHNak0LdOF7lqCSWJtelQW
wHBd2Fd4giuPT1Y6ik7zT4niag5G0BFYxqkDwdtSbQXb+81c7OQjmCtjaet8C7aE
CkfbxsC6xj6050ngM+sliGFfMsP+6sbRUf7uJO2dpA3/R3q5QvqvLfkzUiJpW0+B
IX7v4NJz5OMC15ULqgWIAjTpDaQg6il8/xD8UVtgb6BnWncR+Axr+Kz8oWKoJH0e
lVc0YWa2Ccf20udvl6pqR3nx0dCU3OSVMLAk2hIa3uItqfBvAgQ9S56NKAvIJzPH
nUBXbMPf9jnM2YsCAwEAAaOBlDCBkTAOBgNVHQ8BAf8EBAMCA6gwHQYDVR0lBBYw
FAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQWBBR0z8yLyqIkdyMZI9Sh0/Um
xLyz5DAfBgNVHSMEGDAWgBT/cHRmICNafdJm7uaJ6nVADziO9jAgBgNVHREEGTAX
gg9zZXJ2ZXIuaG9tZS5uZXSHBH8AAAEwDQYJKoZIhvcNAQELBQADggEBACcd7Egz
9Zg7WGRRed1AU2cZ9SYeK7yHdt9AuqUjBP+46DKA1Caw5FHgATtQLGvlgtLkcodo
gfVsPaUDCx3Cwk1PkRP0Rorrve/J2tj3q3lafhszBlXAiKOu66wJxl2X+OAlGZhC
iYwdDzQ3qdyGPQYH37++wA8ysXY1zPj5eBK6Qkzbu7R7t2672RuErtY3SCeiIwUS
NuCTHDPjMsaP89tw7JHZKZltOjjTZ1lwYlaE6T95eSQnvcPnk/GlTSjYwOM7tic4
llDh2bWO1IPL2vizl6RZ5vWf3ajiI6kKKeO+3sqBo7zOFc7pBNdiWNKRSCZuipzY
RzleVBoxWW/hAzI=
-----END CERTIFICATE-----
`
	sn = "6094DEC88A29FC16A180D9DF05893C482E6DFA2D"
)

var _ = Describe("Auth", func() {
	Describe("Registration", func() {
		It("device access registration endpoint with success", func() {
			cert, _ := decodeCertificate(bytes.NewBufferString(registrationCertificate).Bytes())
			service := auth.New(nil, nil)
			newCtx, err := service.Auth(context.TODO(), "/EdgeService/Register", "deviceID", []*x509.Certificate{cert})
			Expect(err).To(BeNil())
			deviceID := newCtx.Value("device_id")
			Expect(deviceID).To(Equal(deviceID))
		})

		It("device access enrol endpoint with success", func() {
			cert, _ := decodeCertificate(bytes.NewBufferString(registrationCertificate).Bytes())
			service := auth.New(nil, nil)
			newCtx, err := service.Auth(context.TODO(), "/EdgeService/Enrol", "deviceID", []*x509.Certificate{cert})
			Expect(err).To(BeNil())
			deviceID := newCtx.Value("device_id")
			Expect(deviceID).To(Equal(deviceID))
		})

		It("access denied when device access other methods with registation certificate", func() {
			cert, _ := decodeCertificate(bytes.NewBufferString(registrationCertificate).Bytes())
			service := auth.New(nil, nil)
			_, err := service.Auth(context.TODO(), "/EdgeService/Admin", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("forbidden to access method"))
		})

		It("access denied when device access enrol endpoint with a real certificate", func() {
			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(nil, nil)
			_, err := service.Auth(context.TODO(), "/EdgeService/Enrol", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("unable to authenticate"))
		})

		It("access denied when device access registration endpoint with a real certificate", func() {
			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(nil, nil)
			_, err := service.Auth(context.TODO(), "/EdgeService/Register", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("unable to authenticate"))
		})
	})

	Describe("Access edge methods", func() {
		It("device access with success edge methods", func() {
			deviceReader := &auth.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:                      "deviceID",
						CertificateSerialNumber: sn,
					}, nil
				},
			}
			certReader := &auth.CertificateReaderMock{
				GetCertificateFunc: func(ctx context.Context, sn string) (entity.CertificateGroup, error) {
					cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
					return entity.CertificateGroup{
						Certificate: cert,
					}, nil
				},
			}

			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(certReader, deviceReader)
			_, err := service.Auth(context.TODO(), "/EdgeService/GetConfiguration", "deviceID", []*x509.Certificate{cert})
			Expect(err).To(BeNil())
		})

		It("device not found", func() {
			deviceReader := &auth.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{}, errors.New("not found")
				},
			}
			certReader := &auth.CertificateReaderMock{
				GetCertificateFunc: func(ctx context.Context, sn string) (entity.CertificateGroup, error) {
					cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
					return entity.CertificateGroup{
						Certificate: cert,
					}, nil
				},
			}

			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(certReader, deviceReader)
			_, err := service.Auth(context.TODO(), "/EdgeService/GetConfiguration", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
		})

		It("unable to get the real certificate", func() {
			deviceReader := &auth.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{CertificateSerialNumber: sn}, nil
				},
			}
			certReader := &auth.CertificateReaderMock{
				GetCertificateFunc: func(ctx context.Context, sn string) (entity.CertificateGroup, error) {
					return entity.CertificateGroup{}, errors.New("not found")
				},
			}

			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(certReader, deviceReader)
			_, err := service.Auth(context.TODO(), "/EdgeService/GetConfiguration", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
		})

		It("device fails to authenticate when serial numbers are different", func() {
			deviceReader := &auth.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:                      "deviceID",
						CertificateSerialNumber: "some serial number",
					}, nil
				},
			}
			certReader := &auth.CertificateReaderMock{
				GetCertificateFunc: func(ctx context.Context, sn string) (entity.CertificateGroup, error) {
					cert, _ := decodeCertificate(bytes.NewBufferString(registrationCertificate).Bytes())
					return entity.CertificateGroup{
						Certificate: cert,
					}, nil
				},
			}

			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(certReader, deviceReader)
			_, err := service.Auth(context.TODO(), "/EdgeService/GetConfiguration", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
		})
		It("device fails to authenticate when the certificate is revoked", func() {
			deviceReader := &auth.DeviceReaderMock{
				GetDeviceFunc: func(ctx context.Context, id string) (entity.Device, error) {
					return entity.Device{
						ID:                      "deviceID",
						CertificateSerialNumber: sn,
					}, nil
				},
			}
			certReader := &auth.CertificateReaderMock{
				GetCertificateFunc: func(ctx context.Context, sn string) (entity.CertificateGroup, error) {
					cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
					return entity.CertificateGroup{
						Certificate: cert,
						IsRevoked:   true,
					}, nil
				},
			}

			cert, _ := decodeCertificate(bytes.NewBufferString(clientCertificate).Bytes())
			service := auth.New(certReader, deviceReader)
			_, err := service.Auth(context.TODO(), "/EdgeService/GetConfiguration", "deviceID", []*x509.Certificate{cert})
			Expect(err).ToNot(BeNil())
		})
	})
})

func decodeCertificate(cert []byte) (*x509.Certificate, error) {
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
