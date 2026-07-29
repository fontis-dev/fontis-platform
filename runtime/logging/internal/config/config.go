package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("LOGGING_SOCKET_PATH", "/run/fontis/logging.sock"),
		DBPath:     getEnv("LOGGING_DB_PATH", "/var/lib/fontis/logging/logging.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
