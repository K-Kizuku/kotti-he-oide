package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config は、アプリケーション設定を保持する構造体
type Config struct {
	// Server設定
	ServerPort string

	// Database設定
	DatabaseURL string

	// Session設定
	SessionTTLMinutes int
}

// Load は、環境変数から設定を読み込む
func Load() (*Config, error) {
	config := &Config{
		ServerPort:        getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		SessionTTLMinutes: getEnvAsInt("SESSION_TTL_MINUTES", 60),
	}

	// 必須の環境変数をチェック
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	return config, nil
}

// getEnv は、環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt は、環境変数を整数として取得し、存在しないまたはパースできない場合はデフォルト値を返す
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
