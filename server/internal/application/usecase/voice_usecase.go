package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/voicevox"
	"github.com/K-Kizuku/kotti-he-oide/pkg/config"
)

// VoiceUseCase は、音声生成のユースケース
type VoiceUseCase struct {
	voicevoxClient *voicevox.Client
	config         *config.Config
}

// NewVoiceUseCase は、新しいVoiceUseCaseを作成する
func NewVoiceUseCase(voicevoxClient *voicevox.Client, cfg *config.Config) *VoiceUseCase {
	return &VoiceUseCase{
		voicevoxClient: voicevoxClient,
		config:         cfg,
	}
}

// GenerateVoice は、テキストから音声を生成し、URLを返す
func (uc *VoiceUseCase) GenerateVoice(text string, speakerID *int) (string, error) {
	// Speaker IDのデフォルト値設定
	actualSpeakerID := uc.config.VoicevoxSpeakerID
	if speakerID != nil {
		actualSpeakerID = *speakerID
	}

	// VOICEVOX APIで音声を生成
	audioData, err := uc.voicevoxClient.GenerateAudio(text, actualSpeakerID)
	if err != nil {
		return "", fmt.Errorf("failed to generate audio from VOICEVOX: %w", err)
	}

	// S3使用の場合とローカル保存の場合で分岐
	if uc.config.UseS3ForAudio {
		// 本番環境: S3にアップロード
		return uc.uploadToS3(audioData, text)
	}

	// 開発環境: ローカルファイルとして保存
	return uc.saveToLocal(audioData, text)
}

// saveToLocal は、音声データをローカルファイルとして保存する
func (uc *VoiceUseCase) saveToLocal(audioData []byte, text string) (string, error) {
	// 出力ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(uc.config.AudioOutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create audio output directory: %w", err)
	}

	// ファイル名を生成（テキストのハッシュ + タイムスタンプ）
	filename := generateFilename(text)
	filepath := filepath.Join(uc.config.AudioOutputDir, filename)

	// ファイルに書き込み
	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to write audio file: %w", err)
	}

	// URLを生成して返す
	audioURL := fmt.Sprintf("%s/%s", uc.config.AudioURLPrefix, filename)
	return audioURL, nil
}

// uploadToS3 は、音声データをS3にアップロードする
// 注意: この実装はスタブです。実際のS3アップロード処理を実装する必要があります。
func (uc *VoiceUseCase) uploadToS3(audioData []byte, text string) (string, error) {
	// TODO: AWS SDK for Goを使ってS3にアップロード
	// 1. AWS SDKのS3クライアントを初期化
	// 2. PutObjectでアップロード
	// 3. S3オブジェクトのURLを返す

	// スタブ実装：エラーを返す
	return "", fmt.Errorf("S3 upload is not implemented yet. Please set USE_S3_FOR_AUDIO=false for development")
}

// generateFilename は、テキストから一意なファイル名を生成する
func generateFilename(text string) string {
	// テキストのSHA256ハッシュを計算
	hash := sha256.Sum256([]byte(text))
	hashStr := hex.EncodeToString(hash[:])

	// タイムスタンプを追加（同じテキストでも時刻が異なれば別ファイルになる）
	timestamp := time.Now().Unix()

	// ファイル名: {hash}_{timestamp}.wav
	return fmt.Sprintf("%s_%d.wav", hashStr[:16], timestamp)
}
