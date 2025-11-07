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

// MySQLPushSubscriptionRepository は、MySQLを使用したPushSubscriptionRepositoryの実装
type MySQLPushSubscriptionRepository struct {
	db *database.DB
}

// NewMySQLPushSubscriptionRepository は、新しいMySQLPushSubscriptionRepositoryを作成する
func NewMySQLPushSubscriptionRepository(db *database.DB) *MySQLPushSubscriptionRepository {
	return &MySQLPushSubscriptionRepository{db: db}
}

// Save は、サブスクリプションを保存する
func (r *MySQLPushSubscriptionRepository) Save(ctx context.Context, subscription *model.PushSubscription) error {
	query := `
		INSERT INTO push_subscriptions (subscription_id, session_id, endpoint, p256dh_key, auth_key, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Conn.ExecContext(
		ctx,
		query,
		subscription.SubscriptionID.String(),
		subscription.SessionID.String(),
		subscription.Endpoint.String(),
		subscription.Keys.P256dh(),
		subscription.Keys.Auth(),
		subscription.IsActive,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save push subscription: %w", err)
	}

	return nil
}

// FindByID は、サブスクリプションIDでサブスクリプションを検索する
func (r *MySQLPushSubscriptionRepository) FindByID(ctx context.Context, subscriptionID valueobject.SubscriptionID) (*model.PushSubscription, error) {
	query := `
		SELECT subscription_id, session_id, endpoint, p256dh_key, auth_key, is_active, created_at, updated_at
		FROM push_subscriptions
		WHERE subscription_id = ?
	`

	var (
		subscription      model.PushSubscription
		subscriptionIDStr string
		sessionIDStr      string
		endpointStr       string
		p256dhKey         string
		authKey           string
	)

	err := r.db.Conn.QueryRowContext(ctx, query, subscriptionID.String()).Scan(
		&subscriptionIDStr,
		&sessionIDStr,
		&endpointStr,
		&p256dhKey,
		&authKey,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewDomainError(
			errors.NOT_FOUND,
			"push subscription not found",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find push subscription: %w", err)
	}

	// Value Objectsを復元
	subID, err := valueobject.NewSubscriptionIDFromString(subscriptionIDStr)
	if err != nil {
		return nil, err
	}
	subscription.SubscriptionID = subID

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	subscription.SessionID = sessionID

	endpoint, err := valueobject.NewPushEndpoint(endpointStr)
	if err != nil {
		return nil, err
	}
	subscription.Endpoint = endpoint

	keys, err := valueobject.NewPushKeys(p256dhKey, authKey)
	if err != nil {
		return nil, err
	}
	subscription.Keys = keys

	return &subscription, nil
}

// FindBySessionID は、セッションIDに紐づくアクティブなサブスクリプションを検索する
func (r *MySQLPushSubscriptionRepository) FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) (*model.PushSubscription, error) {
	query := `
		SELECT subscription_id, session_id, endpoint, p256dh_key, auth_key, is_active, created_at, updated_at
		FROM push_subscriptions
		WHERE session_id = ? AND is_active = TRUE
		ORDER BY created_at DESC
		LIMIT 1
	`

	var (
		subscription      model.PushSubscription
		subscriptionIDStr string
		sessionIDStr      string
		endpointStr       string
		p256dhKey         string
		authKey           string
	)

	err := r.db.Conn.QueryRowContext(ctx, query, sessionID.String()).Scan(
		&subscriptionIDStr,
		&sessionIDStr,
		&endpointStr,
		&p256dhKey,
		&authKey,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewDomainError(
			errors.NOT_FOUND,
			"push subscription not found for session",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find push subscription by session: %w", err)
	}

	// Value Objectsを復元
	subID, err := valueobject.NewSubscriptionIDFromString(subscriptionIDStr)
	if err != nil {
		return nil, err
	}
	subscription.SubscriptionID = subID

	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	subscription.SessionID = sid

	endpoint, err := valueobject.NewPushEndpoint(endpointStr)
	if err != nil {
		return nil, err
	}
	subscription.Endpoint = endpoint

	keys, err := valueobject.NewPushKeys(p256dhKey, authKey)
	if err != nil {
		return nil, err
	}
	subscription.Keys = keys

	return &subscription, nil
}

// Update は、サブスクリプションを更新する
func (r *MySQLPushSubscriptionRepository) Update(ctx context.Context, subscription *model.PushSubscription) error {
	query := `
		UPDATE push_subscriptions
		SET is_active = ?, updated_at = ?
		WHERE subscription_id = ?
	`

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		subscription.IsActive,
		subscription.UpdatedAt,
		subscription.SubscriptionID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update push subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewDomainError(
			errors.NOT_FOUND,
			"push subscription not found",
			nil,
		)
	}

	return nil
}

// Delete は、サブスクリプションを削除する
func (r *MySQLPushSubscriptionRepository) Delete(ctx context.Context, subscriptionID valueobject.SubscriptionID) error {
	query := `DELETE FROM push_subscriptions WHERE subscription_id = ?`

	result, err := r.db.Conn.ExecContext(ctx, query, subscriptionID.String())
	if err != nil {
		return fmt.Errorf("failed to delete push subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewDomainError(
			errors.NOT_FOUND,
			"push subscription not found",
			nil,
		)
	}

	return nil
}
