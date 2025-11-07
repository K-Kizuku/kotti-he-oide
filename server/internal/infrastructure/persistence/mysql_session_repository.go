package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/database"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// MySQLSessionRepository は、MySQLを使用したSessionRepositoryの実装
type MySQLSessionRepository struct {
	db *database.DB
}

// NewMySQLSessionRepository は、新しいMySQLSessionRepositoryを作成する
func NewMySQLSessionRepository(db *database.DB) *MySQLSessionRepository {
	return &MySQLSessionRepository{db: db}
}

// Save は、セッションを保存する
func (r *MySQLSessionRepository) Save(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (session_id, current_scene, s6_started_at, created_at, expires_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Conn.ExecContext(
		ctx,
		query,
		session.SessionID.String(),
		session.CurrentScene,
		session.S6StartedAt,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// FindByID は、セッションIDでセッションを検索する
func (r *MySQLSessionRepository) FindByID(ctx context.Context, sessionID valueobject.SessionID) (*model.Session, error) {
	query := `
		SELECT session_id, current_scene, s6_started_at, created_at, expires_at
		FROM sessions
		WHERE session_id = ?
	`

	var session model.Session
	var sessionIDStr string

	err := r.db.Conn.QueryRowContext(ctx, query, sessionID.String()).Scan(
		&sessionIDStr,
		&session.CurrentScene,
		&session.S6StartedAt,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewDomainError(
			errors.SESSION_NOT_FOUND,
			"session not found",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find session: %w", err)
	}

	// SessionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	session.SessionID = sid

	return &session, nil
}

// Update は、セッションを更新する
func (r *MySQLSessionRepository) Update(ctx context.Context, session *model.Session) error {
	query := `
		UPDATE sessions
		SET current_scene = ?, s6_started_at = ?
		WHERE session_id = ?
	`

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		session.CurrentScene,
		session.S6StartedAt,
		session.SessionID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewDomainError(
			errors.SESSION_NOT_FOUND,
			"session not found",
			nil,
		)
	}

	return nil
}

// Delete は、セッションを削除する
func (r *MySQLSessionRepository) Delete(ctx context.Context, sessionID valueobject.SessionID) error {
	query := `DELETE FROM sessions WHERE session_id = ?`

	result, err := r.db.Conn.ExecContext(ctx, query, sessionID.String())
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewDomainError(
			errors.SESSION_NOT_FOUND,
			"session not found",
			nil,
		)
	}

	return nil
}

// DeleteExpiredSessions は、有効期限切れのセッションを削除する
func (r *MySQLSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`

	_, err := r.db.Conn.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
