package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/database"
)

// MySQLPushLogRepository は、MySQLを使用したPushLogRepositoryの実装
type MySQLPushLogRepository struct {
	db *database.DB
}

// NewMySQLPushLogRepository は、新しいMySQLPushLogRepositoryを作成する
func NewMySQLPushLogRepository(db *database.DB) *MySQLPushLogRepository {
	return &MySQLPushLogRepository{db: db}
}

// Save は、プッシュ通知送信ログを保存する
func (r *MySQLPushLogRepository) Save(ctx context.Context, log *model.PushLog) error {
	query := `
		INSERT INTO push_logs (subscription_id, session_id, title, message, success, status_code, error_message, sent_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// error_messageがNULLの場合
	var errorMessage sql.NullString
	if log.ErrorMessage != "" {
		errorMessage = sql.NullString{String: log.ErrorMessage, Valid: true}
	}

	// titleがNULLの場合
	var title sql.NullString
	if log.Title != "" {
		title = sql.NullString{String: log.Title, Valid: true}
	}

	// status_codeがNULLの場合
	var statusCode sql.NullInt64
	if log.StatusCode != 0 {
		statusCode = sql.NullInt64{Int64: int64(log.StatusCode), Valid: true}
	}

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		log.SubscriptionID.String(),
		log.SessionID.String(),
		title,
		log.Message,
		log.Success,
		statusCode,
		errorMessage,
		log.SentAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save push log: %w", err)
	}

	// IDを取得
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	log.ID = id
	return nil
}

// FindBySessionID は、セッションIDに紐づくログを取得する
func (r *MySQLPushLogRepository) FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) ([]*model.PushLog, error) {
	query := `
		SELECT id, subscription_id, session_id, title, message, success, status_code, error_message, sent_at
		FROM push_logs
		WHERE session_id = ?
		ORDER BY sent_at DESC
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find push logs by session: %w", err)
	}
	defer rows.Close()

	var logs []*model.PushLog

	for rows.Next() {
		var (
			log               model.PushLog
			subscriptionIDStr string
			sessionIDStr      string
			title             sql.NullString
			statusCode        sql.NullInt64
			errorMessage      sql.NullString
		)

		err := rows.Scan(
			&log.ID,
			&subscriptionIDStr,
			&sessionIDStr,
			&title,
			&log.Message,
			&log.Success,
			&statusCode,
			&errorMessage,
			&log.SentAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan push log: %w", err)
		}

		// Value Objectsを復元
		subID, err := valueobject.NewSubscriptionIDFromString(subscriptionIDStr)
		if err != nil {
			return nil, err
		}
		log.SubscriptionID = subID

		sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
		if err != nil {
			return nil, err
		}
		log.SessionID = sid

		// NULL許容フィールドを処理
		if title.Valid {
			log.Title = title.String
		}
		if statusCode.Valid {
			log.StatusCode = int(statusCode.Int64)
		}
		if errorMessage.Valid {
			log.ErrorMessage = errorMessage.String
		}

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating push log rows: %w", err)
	}

	return logs, nil
}

// FindBySubscriptionID は、サブスクリプションIDに紐づくログを取得する
func (r *MySQLPushLogRepository) FindBySubscriptionID(ctx context.Context, subscriptionID valueobject.SubscriptionID) ([]*model.PushLog, error) {
	query := `
		SELECT id, subscription_id, session_id, title, message, success, status_code, error_message, sent_at
		FROM push_logs
		WHERE subscription_id = ?
		ORDER BY sent_at DESC
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, subscriptionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find push logs by subscription: %w", err)
	}
	defer rows.Close()

	var logs []*model.PushLog

	for rows.Next() {
		var (
			log               model.PushLog
			subscriptionIDStr string
			sessionIDStr      string
			title             sql.NullString
			statusCode        sql.NullInt64
			errorMessage      sql.NullString
		)

		err := rows.Scan(
			&log.ID,
			&subscriptionIDStr,
			&sessionIDStr,
			&title,
			&log.Message,
			&log.Success,
			&statusCode,
			&errorMessage,
			&log.SentAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan push log: %w", err)
		}

		// Value Objectsを復元
		subID, err := valueobject.NewSubscriptionIDFromString(subscriptionIDStr)
		if err != nil {
			return nil, err
		}
		log.SubscriptionID = subID

		sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
		if err != nil {
			return nil, err
		}
		log.SessionID = sid

		// NULL許容フィールドを処理
		if title.Valid {
			log.Title = title.String
		}
		if statusCode.Valid {
			log.StatusCode = int(statusCode.Int64)
		}
		if errorMessage.Valid {
			log.ErrorMessage = errorMessage.String
		}

		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating push log rows: %w", err)
	}

	return logs, nil
}
