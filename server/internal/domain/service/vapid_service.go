package service

import (
	"os"

	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// VAPIDService は、VAPID鍵の管理を担当するドメインサービス
type VAPIDService interface {
	// GetPublicKey は、VAPID公開鍵を取得する（クライアント側で使用）
	GetPublicKey() string

	// GetPrivateKey は、VAPID秘密鍵を取得する（サーバー側での署名に使用）
	GetPrivateKey() string
}

// vapidService は、VAPIDServiceの実装
type vapidService struct {
	publicKey  string
	privateKey string
}

// NewVAPIDService は、新しいVAPIDServiceを作成する
// 環境変数からVAPID鍵を読み込む
func NewVAPIDService() (VAPIDService, error) {
	publicKey := os.Getenv("VAPID_PUBLIC_KEY")
	if publicKey == "" {
		return nil, errors.NewDomainError(
			errors.INVALID_INPUT,
			"VAPID_PUBLIC_KEY environment variable is not set",
			nil,
		)
	}

	privateKey := os.Getenv("VAPID_PRIVATE_KEY")
	if privateKey == "" {
		return nil, errors.NewDomainError(
			errors.INVALID_INPUT,
			"VAPID_PRIVATE_KEY environment variable is not set",
			nil,
		)
	}

	return &vapidService{
		publicKey:  publicKey,
		privateKey: privateKey,
	}, nil
}

// GetPublicKey は、VAPID公開鍵を取得する
func (s *vapidService) GetPublicKey() string {
	return s.publicKey
}

// GetPrivateKey は、VAPID秘密鍵を取得する
func (s *vapidService) GetPrivateKey() string {
	return s.privateKey
}
