package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("MARKETPLACE_CLIENT_SOCKET_PATH", "/run/fontis/marketplace-client.sock"),
		DBPath:     getEnv("MARKETPLACE_CLIENT_DB_PATH", "/var/lib/fontis/marketplace-client/marketplace.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
