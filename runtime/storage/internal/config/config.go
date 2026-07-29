package config

import (
	"os"
	"time"
)

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	socketPath := os.Getenv("STORAGE_SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/var/run/fontis/storage.sock"
	}

	dbPath := os.Getenv("STORAGE_DB_PATH")
	if dbPath == "" {
		dbPath = "/var/lib/fontis/storage/storage.db"
	}

	return &Config{
		SocketPath: socketPath,
		DBPath:     dbPath,
	}
}

const ShutdownTimeout = 5 * time.Second
