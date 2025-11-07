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

	// VOICEVOX設定
	VoicevoxAPIURL    string // VOICEVOX API URL（例: http://127.0.0.1:50021）
	VoicevoxSpeakerID int    // 青山龍星(しっとり)のSpeaker ID（デフォルト: 84）
	AudioOutputDir    string // 音声ファイルの出力ディレクトリ
	AudioURLPrefix    string // 音声ファイルのURLプレフィックス
	UseS3ForAudio     bool   // S3を使用するかどうか（本番環境用）
	S3BucketName      string // S3バケット名
	S3Region          string // S3リージョン
}

// Load は、環境変数から設定を読み込む
func Load() (*Config, error) {
	config := &Config{
		ServerPort:        getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		SessionTTLMinutes: getEnvAsInt("SESSION_TTL_MINUTES", 60),

		// VOICEVOX設定
		VoicevoxAPIURL:    getEnv("VOICEVOX_API_URL", "http://127.0.0.1:50021"),
		VoicevoxSpeakerID: getEnvAsInt("VOICEVOX_SPEAKER_ID", 84),
		AudioOutputDir:    getEnv("AUDIO_OUTPUT_DIR", "./audio_files"),
		AudioURLPrefix:    getEnv("AUDIO_URL_PREFIX", "/audio"),
		UseS3ForAudio:     getEnvAsBool("USE_S3_FOR_AUDIO", false),
		S3BucketName:      getEnv("S3_BUCKET_NAME", ""),
		S3Region:          getEnv("S3_REGION", "ap-northeast-1"),
	}

	// 必須の環境変数をチェック
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// S3使用時のバケット名チェック
	if config.UseS3ForAudio && config.S3BucketName == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME is required when USE_S3_FOR_AUDIO is true")
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

// getEnvAsBool は、環境変数をブール値として取得し、存在しないまたはパースできない場合はデフォルト値を返す
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
