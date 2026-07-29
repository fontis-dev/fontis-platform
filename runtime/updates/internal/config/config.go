package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("UPDATES_SOCKET_PATH", "/run/fontis/updates.sock"),
		DBPath:     getEnv("UPDATES_DB_PATH", "/var/lib/fontis/updates/updates.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
