package voicevox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client は、VOICEVOX EngineのAPIクライアント
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient は、新しいVOICEVOXクライアントを作成する
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AudioQuery は、VOICEVOX APIのAudioQueryレスポンス
// 詳細な構造はVOICEVOX APIのスキーマに準拠
type AudioQuery struct {
	AccentPhrases      []AccentPhrase `json:"accent_phrases"`
	SpeedScale         float64        `json:"speedScale"`
	PitchScale         float64        `json:"pitchScale"`
	IntonationScale    float64        `json:"intonationScale"`
	VolumeScale        float64        `json:"volumeScale"`
	PrePhonemeLength   float64        `json:"prePhonemeLength"`
	PostPhonemeLength  float64        `json:"postPhonemeLength"`
	OutputSamplingRate int            `json:"outputSamplingRate"`
	OutputStereo       bool           `json:"outputStereo"`
	Kana               string         `json:"kana"`
}

// AccentPhrase は、アクセント句の情報
type AccentPhrase struct {
	Moras           []Mora `json:"moras"`
	Accent          int    `json:"accent"`
	PauseMora       *Mora  `json:"pause_mora"`
	IsInterrogative bool   `json:"is_interrogative"`
}

// Mora は、モーラ（音素）の情報
type Mora struct {
	Text            string   `json:"text"`
	Consonant       *string  `json:"consonant"`
	ConsonantLength *float64 `json:"consonant_length"`
	Vowel           string   `json:"vowel"`
	VowelLength     float64  `json:"vowel_length"`
	Pitch           float64  `json:"pitch"`
}

// CreateAudioQuery は、テキストからAudioQueryを作成する
// POST /audio_query?text={text}&speaker={speaker}
func (c *Client) CreateAudioQuery(text string, speakerID int) (*AudioQuery, error) {
	// URLエンコード
	params := url.Values{}
	params.Set("text", text)
	params.Set("speaker", strconv.Itoa(speakerID))

	reqURL := fmt.Sprintf("%s/audio_query?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Post(reqURL, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VOICEVOX API returned status %d: %s", resp.StatusCode, string(body))
	}

	var audioQuery AudioQuery
	if err := json.NewDecoder(resp.Body).Decode(&audioQuery); err != nil {
		return nil, fmt.Errorf("failed to decode audio query response: %w", err)
	}

	return &audioQuery, nil
}

// Synthesis は、AudioQueryから音声を合成する
// POST /synthesis?speaker={speaker} + body: AudioQuery
func (c *Client) Synthesis(audioQuery *AudioQuery, speakerID int) ([]byte, error) {
	// AudioQueryをJSONにエンコード
	jsonData, err := json.Marshal(audioQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to encode audio query: %w", err)
	}

	// URLパラメータ
	params := url.Values{}
	params.Set("speaker", strconv.Itoa(speakerID))

	reqURL := fmt.Sprintf("%s/synthesis?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VOICEVOX API returned status %d: %s", resp.StatusCode, string(body))
	}

	// WAVファイルを読み込む
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	return audioData, nil
}

// GenerateAudio は、テキストから音声を生成する（AudioQuery作成 + 音声合成）
func (c *Client) GenerateAudio(text string, speakerID int) ([]byte, error) {
	// Step 1: AudioQueryを作成
	audioQuery, err := c.CreateAudioQuery(text, speakerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio query: %w", err)
	}

	// Step 2: 音声を合成
	audioData, err := c.Synthesis(audioQuery, speakerID)
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize audio: %w", err)
	}

	return audioData, nil
}

// GetSpeakers は、利用可能なスピーカー一覧を取得する
// GET /speakers
func (c *Client) GetSpeakers() ([]Speaker, error) {
	reqURL := fmt.Sprintf("%s/speakers", c.baseURL)

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get speakers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VOICEVOX API returned status %d: %s", resp.StatusCode, string(body))
	}

	var speakers []Speaker
	if err := json.NewDecoder(resp.Body).Decode(&speakers); err != nil {
		return nil, fmt.Errorf("failed to decode speakers response: %w", err)
	}

	return speakers, nil
}

// Speaker は、スピーカー情報
type Speaker struct {
	Name              string                 `json:"name"`
	SpeakerUUID       string                 `json:"speaker_uuid"`
	Styles            []Style                `json:"styles"`
	Version           string                 `json:"version"`
	SupportedFeatures map[string]interface{} `json:"supported_features"`
}

// Style は、スピーカーのスタイル情報
type Style struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
