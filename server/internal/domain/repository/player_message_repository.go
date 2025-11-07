package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// PlayerMessageRepository は、プレイヤーメッセージの永続化を担当するリポジトリインターフェース
type PlayerMessageRepository interface {
	// Save は、メッセージを保存する
	Save(ctx context.Context, message *model.PlayerMessage) error

	// FindByPlaceID は、場所IDでメッセージ一覧を取得する（最新順）
	FindByPlaceID(ctx context.Context, placeID valueobject.PlaceID, limit int) ([]*model.PlayerMessage, error)

	// FindAll は、全メッセージを取得する（最新順、limit付き）
	FindAll(ctx context.Context, limit int) ([]*model.PlayerMessage, error)

	// Count は、メッセージの総数をカウントする
	Count(ctx context.Context) (int, error)
}
