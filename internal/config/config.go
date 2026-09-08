package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds runtime configuration settings for IP-SAKTI Sahayak
type Config struct {
	Port               string
	DatabaseURL        string
	AllowedOrigins     string
	LLMProvider        string
	LLMAPIKey          string
	EmbeddingAPIKey    string
	ZeroDataRetention  bool
	StaticWebDir       string
}

// Load reads configuration from .env and environment variables
func Load() *Config {
	// Attempt loading .env if available, ignore error in production container environments
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] No .env file found, using system environment variables")
	}

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ipsakti?sslmode=disable")
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "*")
	llmProvider := getEnv("LLM_PROVIDER", "mock")
	llmKey := os.Getenv("LLM_API_KEY")
	embeddingKey := os.Getenv("EMBEDDING_API_KEY")
	staticDir := getEnv("STATIC_WEB_DIR", "./web")

	zdrStr := getEnv("ZERO_DATA_RETENTION", "true")
	zdr, err := strconv.ParseBool(zdrStr)
	if err != nil {
		zdr = true
	}

	return &Config{
		Port:              port,
		DatabaseURL:       dbURL,
		AllowedOrigins:    allowedOrigins,
		LLMProvider:       llmProvider,
		LLMAPIKey:         llmKey,
		EmbeddingAPIKey:   embeddingKey,
		ZeroDataRetention: zdr,
		StaticWebDir:      staticDir,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
