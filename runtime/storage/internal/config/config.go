package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("STORAGE_SOCKET_PATH", "/run/fontis/storage.sock"),
		DBPath:     getEnv("STORAGE_DB_PATH", "/var/lib/fontis/storage/storage.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
