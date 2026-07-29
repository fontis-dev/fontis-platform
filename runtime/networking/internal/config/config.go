package config

import (
	"os"
	"time"
)

type Config struct {
	SocketPath string
}

func Load() *Config {
	socketPath := os.Getenv("NETWORKING_SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/var/run/fontis/networking.sock"
	}

	return &Config{
		SocketPath: socketPath,
	}
}

const ShutdownTimeout = 5 * time.Second
