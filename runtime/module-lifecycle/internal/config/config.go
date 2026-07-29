package config

import "os"

type Config struct {
	SocketPath string
	DBPath     string
}

func Load() *Config {
	return &Config{
		SocketPath: getEnv("MODULE_LIFECYCLE_SOCKET_PATH", "/run/fontis/module-lifecycle.sock"),
		DBPath:     getEnv("MODULE_LIFECYCLE_DB_PATH", "/var/lib/fontis/module-lifecycle/modules.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
