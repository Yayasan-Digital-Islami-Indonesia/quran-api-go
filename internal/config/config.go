package config

import "os"

type Config struct {
	DBPath         string
	ServerHost     string
	ServerPort     string
	AllowedOrigins string
	AppVersion     string
	LogLevel       string
}

func Load() Config {
	cfg := Config{
		DBPath:         getenv("DB_PATH", "./data/quran.db"),
		ServerHost:     getenv("SERVER_HOST", "0.0.0.0"),
		ServerPort:     getenv("SERVER_PORT", "8080"),
		AllowedOrigins: getenv("ALLOWED_ORIGINS", ""),
		AppVersion:     getenv("APP_VERSION", "1.0.0"),
		LogLevel:       getenv("LOG_LEVEL", "info"),
	}

	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
