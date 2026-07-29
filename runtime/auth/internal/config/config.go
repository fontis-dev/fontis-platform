package config

import (
	"os"
	"time"
)

type Config struct {
	SocketPath  string
	DBPath      string
	TLSCertPath string
	TLSKeyPath  string
	TLSCAPath   string
	DevTLS      bool
}

func Load() *Config {
	socketPath := os.Getenv("AUTH_SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/var/run/fontis/auth.sock"
	}

	dbPath := os.Getenv("AUTH_DB_PATH")
	if dbPath == "" {
		dbPath = "/var/lib/fontis/auth/auth.db"
	}

	tlsCertPath := os.Getenv("AUTH_TLS_CERT")
	tlsKeyPath := os.Getenv("AUTH_TLS_KEY")
	tlsCAPath := os.Getenv("AUTH_TLS_CA")

	devTLS := os.Getenv("FONTIS_DEV") == "1" || (tlsCertPath == "" && tlsKeyPath == "")

	return &Config{
		SocketPath:  socketPath,
		DBPath:      dbPath,
		TLSCertPath: tlsCertPath,
		TLSKeyPath:  tlsKeyPath,
		TLSCAPath:   tlsCAPath,
		DevTLS:      devTLS,
	}
}

const ShutdownTimeout = 5 * time.Second
