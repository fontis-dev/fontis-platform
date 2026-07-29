package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("IDENTITY_SOCKET_PATH", "/run/fontis/identity.sock"),
		DBPath:     getEnv("IDENTITY_DB_PATH", "/var/lib/fontis/identity/identity.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
