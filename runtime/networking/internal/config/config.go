package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("NETWORKING_SOCKET_PATH", "/run/fontis/networking.sock"),
		DBPath:     getEnv("NETWORKING_DB_PATH", "/var/lib/fontis/networking/networking.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
