package config

import "os"

type Config struct {
	DSN       string
	JWTSecret string
	Port      string
}

func Load() *Config {
	return &Config{
		DSN:       os.Getenv("BLOG_MYSQL_DSN"),
		JWTSecret: getEnv("BLOG_JWT_SECRET", "default_jwt_secret"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
