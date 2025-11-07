package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// SessionRepository は、セッションの永続化を担当するリポジトリインターフェース
type SessionRepository interface {
	// Save は、セッションを保存する
	Save(ctx context.Context, session *model.Session) error

	// FindByID は、セッションIDでセッションを検索する
	FindByID(ctx context.Context, sessionID valueobject.SessionID) (*model.Session, error)

	// Update は、セッションを更新する
	Update(ctx context.Context, session *model.Session) error

	// Delete は、セッションを削除する
	Delete(ctx context.Context, sessionID valueobject.SessionID) error

	// DeleteExpiredSessions は、有効期限切れのセッションを削除する
	DeleteExpiredSessions(ctx context.Context) error
}
