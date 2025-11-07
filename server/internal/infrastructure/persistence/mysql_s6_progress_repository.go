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

// MySQLS6ProgressRepository は、MySQLを使用したS6ProgressRepositoryの実装
type MySQLS6ProgressRepository struct {
	db *database.DB
}

// NewMySQLS6ProgressRepository は、新しいMySQLS6ProgressRepositoryを作成する
func NewMySQLS6ProgressRepository(db *database.DB) *MySQLS6ProgressRepository {
	return &MySQLS6ProgressRepository{db: db}
}

// Save は、進捗を保存する
func (r *MySQLS6ProgressRepository) Save(ctx context.Context, progress *model.S6Progress) error {
	query := `
		INSERT INTO session_s6_progress (session_id, place_id, verified, verified_by, quiz_id, answered, correct, verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			verified = VALUES(verified),
			verified_by = VALUES(verified_by),
			quiz_id = VALUES(quiz_id),
			answered = VALUES(answered),
			correct = VALUES(correct),
			verified_at = VALUES(verified_at)
	`

	var quizIDVal interface{}
	if progress.QuizID != nil {
		quizIDVal = progress.QuizID.String()
	} else {
		quizIDVal = nil
	}

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		progress.SessionID.String(),
		progress.PlaceID.String(),
		progress.Verified,
		progress.VerifiedBy,
		quizIDVal,
		progress.Answered,
		progress.Correct,
		progress.VerifiedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save s6 progress: %w", err)
	}

	if progress.ID == 0 {
		id, err := result.LastInsertId()
		if err == nil {
			progress.ID = int(id)
		}
	}

	return nil
}

// FindBySessionID は、セッションIDで全進捗を取得する
func (r *MySQLS6ProgressRepository) FindBySessionID(
	ctx context.Context,
	sessionID valueobject.SessionID,
) ([]*model.S6Progress, error) {
	query := `
		SELECT id, session_id, place_id, verified, verified_by, quiz_id, answered, correct, verified_at
		FROM session_s6_progress
		WHERE session_id = ?
		ORDER BY place_id ASC
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find s6 progress: %w", err)
	}
	defer rows.Close()

	var progressList []*model.S6Progress
	for rows.Next() {
		progress, err := r.scanS6Progress(rows)
		if err != nil {
			return nil, err
		}
		progressList = append(progressList, progress)
	}

	return progressList, nil
}

// FindBySessionIDAndPlaceID は、セッションIDと場所IDで進捗を取得する
func (r *MySQLS6ProgressRepository) FindBySessionIDAndPlaceID(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
) (*model.S6Progress, error) {
	query := `
		SELECT id, session_id, place_id, verified, verified_by, quiz_id, answered, correct, verified_at
		FROM session_s6_progress
		WHERE session_id = ? AND place_id = ?
	`

	progress := &model.S6Progress{}
	var sessionIDStr string
	var placeIDStr string
	var quizIDStr sql.NullString

	err := r.db.Conn.QueryRowContext(ctx, query, sessionID.String(), placeID.String()).Scan(
		&progress.ID,
		&sessionIDStr,
		&placeIDStr,
		&progress.Verified,
		&progress.VerifiedBy,
		&quizIDStr,
		&progress.Answered,
		&progress.Correct,
		&progress.VerifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewDomainError(
			errors.NOT_FOUND,
			"s6 progress not found",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find s6 progress: %w", err)
	}

	// SessionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	progress.SessionID = sid

	// PlaceIDを復元
	pid, err := valueobject.NewPlaceID(placeIDStr)
	if err != nil {
		return nil, err
	}
	progress.PlaceID = pid

	// QuizIDを復元
	if quizIDStr.Valid {
		qid, err := valueobject.NewQuizIDFromString(quizIDStr.String)
		if err != nil {
			return nil, err
		}
		progress.QuizID = &qid
	}

	return progress, nil
}

// Update は、進捗を更新する
func (r *MySQLS6ProgressRepository) Update(ctx context.Context, progress *model.S6Progress) error {
	query := `
		UPDATE session_s6_progress
		SET verified = ?, verified_by = ?, quiz_id = ?, answered = ?, correct = ?, verified_at = ?
		WHERE session_id = ? AND place_id = ?
	`

	var quizIDVal interface{}
	if progress.QuizID != nil {
		quizIDVal = progress.QuizID.String()
	} else {
		quizIDVal = nil
	}

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		progress.Verified,
		progress.VerifiedBy,
		quizIDVal,
		progress.Answered,
		progress.Correct,
		progress.VerifiedAt,
		progress.SessionID.String(),
		progress.PlaceID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update s6 progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.NewDomainError(
			errors.NOT_FOUND,
			"s6 progress not found",
			nil,
		)
	}

	return nil
}

// CountCompleted は、完了した場所の数をカウントする
func (r *MySQLS6ProgressRepository) CountCompleted(ctx context.Context, sessionID valueobject.SessionID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM session_s6_progress
		WHERE session_id = ? AND verified = true AND answered = true AND correct = true
	`

	var count int
	err := r.db.Conn.QueryRowContext(ctx, query, sessionID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count completed s6 progress: %w", err)
	}

	return count, nil
}

// scanS6Progress は、rowsからS6Progressをスキャンする
func (r *MySQLS6ProgressRepository) scanS6Progress(rows *sql.Rows) (*model.S6Progress, error) {
	var progress model.S6Progress
	var sessionIDStr string
	var placeIDStr string
	var quizIDStr sql.NullString

	err := rows.Scan(
		&progress.ID,
		&sessionIDStr,
		&placeIDStr,
		&progress.Verified,
		&progress.VerifiedBy,
		&quizIDStr,
		&progress.Answered,
		&progress.Correct,
		&progress.VerifiedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan s6 progress: %w", err)
	}

	// SessionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	progress.SessionID = sid

	// PlaceIDを復元
	pid, err := valueobject.NewPlaceID(placeIDStr)
	if err != nil {
		return nil, err
	}
	progress.PlaceID = pid

	// QuizIDを復元
	if quizIDStr.Valid {
		qid, err := valueobject.NewQuizIDFromString(quizIDStr.String)
		if err != nil {
			return nil, err
		}
		progress.QuizID = &qid
	}

	return &progress, nil
}
