package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// S6ProgressRepository は、S6進捗の永続化を担当するリポジトリインターフェース
type S6ProgressRepository interface {
	// Save は、進捗を保存する
	Save(ctx context.Context, progress *model.S6Progress) error

	// FindBySessionID は、セッションIDで全進捗を取得する
	FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) ([]*model.S6Progress, error)

	// FindBySessionIDAndPlaceID は、セッションIDと場所IDで進捗を取得する
	FindBySessionIDAndPlaceID(
		ctx context.Context,
		sessionID valueobject.SessionID,
		placeID valueobject.PlaceID,
	) (*model.S6Progress, error)

	// Update は、進捗を更新する
	Update(ctx context.Context, progress *model.S6Progress) error

	// CountCompleted は、完了した場所の数をカウントする
	CountCompleted(ctx context.Context, sessionID valueobject.SessionID) (int, error)
}
