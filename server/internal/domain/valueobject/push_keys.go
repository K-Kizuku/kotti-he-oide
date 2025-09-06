package valueobject

import (
	"encoding/base64"
	"fmt"
)

type P256dhKey struct {
	value string
}

func NewP256dhKey(key string) (P256dhKey, error) {
	if key == "" {
		return P256dhKey{}, fmt.Errorf("p256dh key cannot be empty")
	}

	// Validate base64url encoding
	_, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return P256dhKey{}, fmt.Errorf("p256dh key must be valid base64url: %w", err)
	}

	return P256dhKey{value: key}, nil
}

func (k P256dhKey) Value() string {
	return k.value
}

func (k P256dhKey) String() string {
	return k.value
}

func (k P256dhKey) Equals(other P256dhKey) bool {
	return k.value == other.value
}

type AuthKey struct {
	value string
}

func NewAuthKey(key string) (AuthKey, error) {
	if key == "" {
		return AuthKey{}, fmt.Errorf("auth key cannot be empty")
	}

	// Validate base64url encoding
	_, err := base64.RawURLEncoding.DecodeString(key)
	if err != nil {
		return AuthKey{}, fmt.Errorf("auth key must be valid base64url: %w", err)
	}

	return AuthKey{value: key}, nil
}

func (k AuthKey) Value() string {
	return k.value
}

func (k AuthKey) String() string {
	return k.value
}

func (k AuthKey) Equals(other AuthKey) bool {
	return k.value == other.value
}

type PushKeys struct {
	p256dh P256dhKey
	auth   AuthKey
}

func NewPushKeys(p256dh P256dhKey, auth AuthKey) PushKeys {
	return PushKeys{
		p256dh: p256dh,
		auth:   auth,
	}
}

func (k PushKeys) P256dh() P256dhKey {
	return k.p256dh
}

func (k PushKeys) Auth() AuthKey {
	return k.auth
}

func (k PushKeys) Equals(other PushKeys) bool {
	return k.p256dh.Equals(other.p256dh) && k.auth.Equals(other.auth)
}