package valueobject

import (
	"encoding/base64"
	"strings"

	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// PushKeys は、Web Push APIで使用されるP256dhキーとAuthキーを表すValue Object
type PushKeys struct {
	p256dh string // 公開鍵（Base64エンコード）
	auth   string // 認証シークレット（Base64エンコード）
}

// NewPushKeys は、新しいPushKeysを作成する
func NewPushKeys(p256dh, auth string) (PushKeys, error) {
	// 空文字列チェック
	if strings.TrimSpace(p256dh) == "" {
		return PushKeys{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"p256dh key cannot be empty",
			nil,
		)
	}
	if strings.TrimSpace(auth) == "" {
		return PushKeys{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"auth key cannot be empty",
			nil,
		)
	}

	// Base64エンコードの検証
	if _, err := base64.StdEncoding.DecodeString(p256dh); err != nil {
		// URL-safe Base64も試す
		if _, err := base64.RawURLEncoding.DecodeString(p256dh); err != nil {
			return PushKeys{}, errors.NewDomainError(
				errors.INVALID_INPUT,
				"p256dh key must be valid Base64",
				err,
			)
		}
	}

	if _, err := base64.StdEncoding.DecodeString(auth); err != nil {
		// URL-safe Base64も試す
		if _, err := base64.RawURLEncoding.DecodeString(auth); err != nil {
			return PushKeys{}, errors.NewDomainError(
				errors.INVALID_INPUT,
				"auth key must be valid Base64",
				err,
			)
		}
	}

	return PushKeys{
		p256dh: p256dh,
		auth:   auth,
	}, nil
}

// P256dh は、P256dh公開鍵を返す
func (k PushKeys) P256dh() string {
	return k.p256dh
}

// Auth は、認証シークレットを返す
func (k PushKeys) Auth() string {
	return k.auth
}

// Equals は、2つのPushKeysが等しいかどうかを判定する
func (k PushKeys) Equals(other PushKeys) bool {
	return k.p256dh == other.p256dh && k.auth == other.auth
}
