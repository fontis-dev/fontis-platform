package fontistls

import (
	"crypto/tls"
	"fmt"
	"log"
)

func DevTLSConfig(serviceName string) (*tls.Config, *tls.Config, error) {
	log.Printf("[tls] generating ephemeral CA and certificate for %s (dev mode)", serviceName)

	ca, err := GenerateCA()
	if err != nil {
		return nil, nil, fmt.Errorf("dev tls generate ca: %w", err)
	}

	sc, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, serviceName)
	if err != nil {
		return nil, nil, fmt.Errorf("dev tls generate cert: %w", err)
	}

	serverCfg, err := NewServerTLSConfig(ca.CertPEM, sc.CertPEM, sc.KeyPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("dev tls server config: %w", err)
	}

	clientCfg, err := NewClientTLSConfig(ca.CertPEM, sc.CertPEM, sc.KeyPEM, serviceName)
	if err != nil {
		return nil, nil, fmt.Errorf("dev tls client config: %w", err)
	}

	return serverCfg, clientCfg, nil
}
