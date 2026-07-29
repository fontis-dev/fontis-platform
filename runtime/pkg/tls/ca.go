package fontistls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

type CA struct {
	CertPEM []byte
	KeyPEM  []byte
}

func GenerateCA() (*CA, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate ca key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("generate ca serial: %w", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "Fontis Device CA",
			Organization: []string{"Fontis"},
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("create ca certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return &CA{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}

func (ca *CA) Certificate() (*x509.Certificate, error) {
	block, _ := pem.Decode(ca.CertPEM)
	if block == nil {
		return nil, fmt.Errorf("decode ca cert pem")
	}
	return x509.ParseCertificate(block.Bytes)
}

func (ca *CA) PrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(ca.KeyPEM)
	if block == nil {
		return nil, fmt.Errorf("decode ca key pem")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
