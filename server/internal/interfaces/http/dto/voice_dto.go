package dto

// GenerateVoiceRequest は、音声生成リクエスト
type GenerateVoiceRequest struct {
	Text      string `json:"text"`       // 音声化するテキスト
	SpeakerID *int   `json:"speaker_id"` // Speaker ID（省略時はデフォルト値を使用）
}

// GenerateVoiceResponse は、音声生成レスポンス
type GenerateVoiceResponse struct {
	AudioURL string `json:"audio_url"` // 生成された音声ファイルのURL
}

// ErrorResponse は、エラーレスポンス
type ErrorResponse struct {
	Error   string `json:"error"`             // エラーコード
	Message string `json:"message"`           // エラーメッセージ
	Details string `json:"details,omitempty"` // 詳細情報（オプション）
}
