package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("BACKUP_SOCKET_PATH", "/run/fontis/backup.sock"),
		DBPath:     getEnv("BACKUP_DB_PATH", "/var/lib/fontis/backup/backup.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
