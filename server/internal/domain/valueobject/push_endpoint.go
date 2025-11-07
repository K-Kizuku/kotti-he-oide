package valueobject

import (
	"net/url"
	"strings"

	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// PushEndpoint は、プッシュ通知のエンドポイントURLを表すValue Object
type PushEndpoint struct {
	value string
}

// NewPushEndpoint は、新しいPushEndpointを作成する
func NewPushEndpoint(endpoint string) (PushEndpoint, error) {
	// 空文字列チェック
	if strings.TrimSpace(endpoint) == "" {
		return PushEndpoint{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"push endpoint cannot be empty",
			nil,
		)
	}

	// URLフォーマットの検証
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return PushEndpoint{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"invalid push endpoint URL format",
			err,
		)
	}

	// HTTPSスキームの検証（開発環境ではHTTPも許可）
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return PushEndpoint{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"push endpoint must use HTTP or HTTPS scheme",
			nil,
		)
	}

	return PushEndpoint{value: endpoint}, nil
}

// String は、PushEndpointを文字列として返す
func (p PushEndpoint) String() string {
	return p.value
}

// Equals は、2つのPushEndpointが等しいかどうかを判定する
func (p PushEndpoint) Equals(other PushEndpoint) bool {
	return p.value == other.value
}
