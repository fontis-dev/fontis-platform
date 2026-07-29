package config

import "os"

type Config struct {
	SocketPath         string
	DBPath             string
	SessionDurationMin int
}

func Load() *Config {
	return &Config{
		SocketPath:         getEnv("AUTH_SOCKET_PATH", "/run/fontis/auth.sock"),
		DBPath:             getEnv("AUTH_DB_PATH", "/var/lib/fontis/auth/auth.db"),
		SessionDurationMin: 15,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
