package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/database"
)

// MySQLPlayerMessageRepository は、MySQLを使用したPlayerMessageRepositoryの実装
type MySQLPlayerMessageRepository struct {
	db *database.DB
}

// NewMySQLPlayerMessageRepository は、新しいMySQLPlayerMessageRepositoryを作成する
func NewMySQLPlayerMessageRepository(db *database.DB) *MySQLPlayerMessageRepository {
	return &MySQLPlayerMessageRepository{db: db}
}

// Save は、メッセージを保存する
func (r *MySQLPlayerMessageRepository) Save(ctx context.Context, message *model.PlayerMessage) error {
	query := `
		INSERT INTO player_messages (session_id, place_id, message_text, created_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		message.SessionID.String(),
		message.PlaceID.String(),
		message.MessageText,
		message.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save player message: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		message.ID = int(id)
	}

	return nil
}

// FindByPlaceID は、場所IDでメッセージ一覧を取得する（最新順）
func (r *MySQLPlayerMessageRepository) FindByPlaceID(
	ctx context.Context,
	placeID valueobject.PlaceID,
	limit int,
) ([]*model.PlayerMessage, error) {
	query := `
		SELECT id, session_id, place_id, message_text, created_at
		FROM player_messages
		WHERE place_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, placeID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find player messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.PlayerMessage
	for rows.Next() {
		message, err := r.scanPlayerMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// FindAll は、全メッセージを取得する（最新順、limit付き）
func (r *MySQLPlayerMessageRepository) FindAll(ctx context.Context, limit int) ([]*model.PlayerMessage, error) {
	query := `
		SELECT id, session_id, place_id, message_text, created_at
		FROM player_messages
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find player messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.PlayerMessage
	for rows.Next() {
		message, err := r.scanPlayerMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// Count は、メッセージの総数をカウントする
func (r *MySQLPlayerMessageRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM player_messages`

	var count int
	err := r.db.Conn.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count player messages: %w", err)
	}

	return count, nil
}

// scanPlayerMessage は、rowsからPlayerMessageをスキャンする
func (r *MySQLPlayerMessageRepository) scanPlayerMessage(rows *sql.Rows) (*model.PlayerMessage, error) {
	var message model.PlayerMessage
	var sessionIDStr string
	var placeIDStr string

	err := rows.Scan(
		&message.ID,
		&sessionIDStr,
		&placeIDStr,
		&message.MessageText,
		&message.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan player message: %w", err)
	}

	// SessionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	message.SessionID = sid

	// PlaceIDを復元
	pid, err := valueobject.NewPlaceID(placeIDStr)
	if err != nil {
		return nil, err
	}
	message.PlaceID = pid

	return &message, nil
}
