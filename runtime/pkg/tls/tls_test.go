package fontistls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestGenerateCA(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}
	if len(ca.CertPEM) == 0 {
		t.Error("expected non-empty cert PEM")
	}
	if len(ca.KeyPEM) == 0 {
		t.Error("expected non-empty key PEM")
	}

	cert, err := ca.Certificate()
	if err != nil {
		t.Fatalf("Certificate: %v", err)
	}
	if !cert.IsCA {
		t.Error("expected CA certificate")
	}
	if cert.Subject.CommonName != "Fontis Device CA" {
		t.Errorf("got CN %q, want %q", cert.Subject.CommonName, "Fontis Device CA")
	}
}

func TestGenerateServiceCert(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}

	sc, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, "test-service")
	if err != nil {
		t.Fatalf("GenerateServiceCert: %v", err)
	}
	if len(sc.CertPEM) == 0 {
		t.Error("expected non-empty cert PEM")
	}
	if len(sc.KeyPEM) == 0 {
		t.Error("expected non-empty key PEM")
	}

	block, _ := pem.Decode(sc.CertPEM)
	if block == nil {
		t.Fatal("decode service cert pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("ParseCertificate: %v", err)
	}
	if cert.Subject.CommonName != "test-service" {
		t.Errorf("got CN %q, want %q", cert.Subject.CommonName, "test-service")
	}

	caCert, err := ca.Certificate()
	if err != nil {
		t.Fatalf("CA Certificate: %v", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(caCert)
	opts := x509.VerifyOptions{Roots: roots}
	if _, err := cert.Verify(opts); err != nil {
		t.Errorf("cert verify against CA: %v", err)
	}
}

func TestServerTLSConfig(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}
	sc, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, "server")
	if err != nil {
		t.Fatalf("GenerateServiceCert: %v", err)
	}

	tlsCfg, err := NewServerTLSConfig(ca.CertPEM, sc.CertPEM, sc.KeyPEM)
	if err != nil {
		t.Fatalf("NewServerTLSConfig: %v", err)
	}

	if len(tlsCfg.Certificates) != 1 {
		t.Error("expected 1 server certificate")
	}
	if tlsCfg.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Error("expected RequireAndVerifyClientCert")
	}
	if tlsCfg.MinVersion != tls.VersionTLS12 {
		t.Error("expected TLS 1.2 minimum")
	}
}

func TestClientTLSConfig(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}
	sc, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, "client")
	if err != nil {
		t.Fatalf("GenerateServiceCert: %v", err)
	}

	tlsCfg, err := NewClientTLSConfig(ca.CertPEM, sc.CertPEM, sc.KeyPEM, "fontis-device")
	if err != nil {
		t.Fatalf("NewClientTLSConfig: %v", err)
	}

	if len(tlsCfg.Certificates) != 1 {
		t.Error("expected 1 client certificate")
	}
	if tlsCfg.ServerName != "fontis-device" {
		t.Errorf("got ServerName %q, want %q", tlsCfg.ServerName, "fontis-device")
	}
	if tlsCfg.MinVersion != tls.VersionTLS12 {
		t.Error("expected TLS 1.2 minimum")
	}
}

func TestMutualTLSCertificateVerification(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}

	serverCert, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, "identity")
	if err != nil {
		t.Fatalf("server GenerateServiceCert: %v", err)
	}
	clientCert, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, "auth")
	if err != nil {
		t.Fatalf("client GenerateServiceCert: %v", err)
	}

	serverTLS, err := NewServerTLSConfig(ca.CertPEM, serverCert.CertPEM, serverCert.KeyPEM)
	if err != nil {
		t.Fatalf("NewServerTLSConfig: %v", err)
	}
	clientTLS, err := NewClientTLSConfig(ca.CertPEM, clientCert.CertPEM, clientCert.KeyPEM, "identity")
	if err != nil {
		t.Fatalf("NewClientTLSConfig: %v", err)
	}

	ln, err := tls.Listen("tcp", "127.0.0.1:0", serverTLS)
	if err != nil {
		t.Fatalf("tls.Listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().String()

	errCh := make(chan error, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer conn.Close()

		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			errCh <- fmt.Errorf("not a tls connection")
			return
		}
		if err := tlsConn.Handshake(); err != nil {
			errCh <- err
			return
		}

		state := tlsConn.ConnectionState()
		if len(state.PeerCertificates) == 0 {
			errCh <- fmt.Errorf("no peer certificates")
			return
		}
		errCh <- nil
	}()

	conn, err := tls.Dial("tcp", addr, clientTLS)
	if err != nil {
		t.Fatalf("tls.Dial: %v", err)
	}
	conn.Close()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server handshake: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for server handshake")
	}
}

func TestRejectClientWithoutCert(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}
	serverCert, err := GenerateServiceCert(ca.CertPEM, ca.KeyPEM, "identity")
	if err != nil {
		t.Fatalf("GenerateServiceCert: %v", err)
	}

	serverTLS, err := NewServerTLSConfig(ca.CertPEM, serverCert.CertPEM, serverCert.KeyPEM)
	if err != nil {
		t.Fatalf("NewServerTLSConfig: %v", err)
	}

	rawLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	defer rawLn.Close()

	addr := rawLn.Addr().String()

	errCh := make(chan error, 1)
	go func() {
		rawConn, err := rawLn.Accept()
		if err != nil {
			errCh <- err
			return
		}
		tlsConn := tls.Server(rawConn, serverTLS)
		if err := tlsConn.Handshake(); err != nil {
			errCh <- nil
			return
		}
		errCh <- fmt.Errorf("expected handshake to fail")
	}()

	clientTLS := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, clientTLS)
	if err == nil {
		conn.SetReadDeadline(time.Now().Add(time.Second))
		_, err = conn.Read([]byte{0})
		conn.Close()
	}
	if err == nil {
		t.Fatal("expected connection to fail without client cert")
	}

	select {
	case hsErr := <-errCh:
		if hsErr != nil {
			t.Fatalf("server handshake: %v", hsErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for server handshake")
	}
}

func TestCAKeyPair(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("GenerateCA: %v", err)
	}

	cert, err := ca.Certificate()
	if err != nil {
		t.Fatalf("Certificate: %v", err)
	}
	if cert.PublicKeyAlgorithm != x509.RSA {
		t.Errorf("got key algo %v, want RSA", cert.PublicKeyAlgorithm)
	}

	key, err := ca.PrivateKey()
	if err != nil {
		t.Fatalf("PrivateKey: %v", err)
	}
	if key.N.BitLen() != 2048 {
		t.Errorf("got key size %d, want 2048", key.N.BitLen())
	}
}

func TestServerConfigValidation(t *testing.T) {
	_, err := NewServerTLSConfig(nil, nil, nil)
	if err == nil {
		t.Error("expected error with nil PEMs")
	}
}

func TestClientConfigValidation(t *testing.T) {
	_, err := NewClientTLSConfig(nil, nil, nil, "server")
	if err == nil {
		t.Error("expected error with nil PEMs")
	}
}


