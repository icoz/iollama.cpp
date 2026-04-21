package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ModelPath   string
	LibPath     string
	Host        string
	Port        int
	MaxTokens   int
	Temperature float64
	ContextSize int
}

// Load загружает конфигурацию из переменных окружения с префиксом IOLLAMA_
func Load() (*Config, error) {
	cfg := &Config{
		ModelPath:   os.Getenv("IOLLAMA_MODEL_PATH"),
		LibPath:     os.Getenv("IOLLAMA_LIB_PATH"),
		Host:        getEnv("IOLLAMA_HOST", "127.0.0.1"),
		Port:        getEnvInt("IOLLAMA_PORT", 8080),
		MaxTokens:   getEnvInt("IOLLAMA_MAX_TOKENS", 2048),
		Temperature: getEnvFloat("IOLLAMA_TEMPERATURE", 0.7),
		ContextSize: getEnvInt("IOLLAMA_CONTEXT_SIZE", 4096),
	}

	if cfg.ModelPath == "" {
		return nil, fmt.Errorf("IOLLAMA_MODEL_PATH is required")
	}
	if _, err := os.Stat(cfg.ModelPath); err != nil {
		return nil, fmt.Errorf("model path does not exist: %w", err)
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}
