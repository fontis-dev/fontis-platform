package fontistls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

func NewServerTLSConfig(caPEM, certPEM, keyPEM []byte) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("append ca cert to pool")
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("load server key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ServerName:   "fontis-device",
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func NewClientTLSConfig(caPEM, certPEM, keyPEM []byte, serverName string) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("append ca cert to pool")
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("load client key pair: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   serverName,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

const DefaultServerName = "fontis-device"
