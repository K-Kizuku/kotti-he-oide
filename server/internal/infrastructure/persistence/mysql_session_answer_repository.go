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

// MySQLSessionAnswerRepository は、MySQLを使用したSessionAnswerRepositoryの実装
type MySQLSessionAnswerRepository struct {
	db *database.DB
}

// NewMySQLSessionAnswerRepository は、新しいMySQLSessionAnswerRepositoryを作成する
func NewMySQLSessionAnswerRepository(db *database.DB) *MySQLSessionAnswerRepository {
	return &MySQLSessionAnswerRepository{db: db}
}

// Save は、回答を保存する（既存の回答がある場合は上書き）
func (r *MySQLSessionAnswerRepository) Save(ctx context.Context, answer *model.SessionAnswer) error {
	query := `
		INSERT INTO session_answers (session_id, question_id, answer_text, answered_at)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE answer_text = VALUES(answer_text), answered_at = VALUES(answered_at)
	`

	result, err := r.db.Conn.ExecContext(
		ctx,
		query,
		answer.SessionID.String(),
		answer.QuestionID.Int(),
		answer.AnswerText,
		answer.AnsweredAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save session answer: %w", err)
	}

	if answer.ID == 0 {
		id, err := result.LastInsertId()
		if err == nil {
			answer.ID = int(id)
		}
	}

	return nil
}

// FindBySessionID は、セッションIDで全回答を取得する
func (r *MySQLSessionAnswerRepository) FindBySessionID(
	ctx context.Context,
	sessionID valueobject.SessionID,
) ([]*model.SessionAnswer, error) {
	query := `
		SELECT id, session_id, question_id, answer_text, answered_at
		FROM session_answers
		WHERE session_id = ?
		ORDER BY question_id ASC
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find session answers: %w", err)
	}
	defer rows.Close()

	var answers []*model.SessionAnswer
	for rows.Next() {
		var answer model.SessionAnswer
		var sessionIDStr string
		var questionIDInt int

		err := rows.Scan(
			&answer.ID,
			&sessionIDStr,
			&questionIDInt,
			&answer.AnswerText,
			&answer.AnsweredAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session answer: %w", err)
		}

		// SessionIDとQuestionIDを復元
		sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
		if err != nil {
			return nil, err
		}
		answer.SessionID = sid

		qid, err := valueobject.NewQuestionID(questionIDInt)
		if err != nil {
			return nil, err
		}
		answer.QuestionID = qid

		answers = append(answers, &answer)
	}

	return answers, nil
}

// FindBySessionIDAndQuestionID は、セッションIDと質問IDで回答を取得する
func (r *MySQLSessionAnswerRepository) FindBySessionIDAndQuestionID(
	ctx context.Context,
	sessionID valueobject.SessionID,
	questionID valueobject.QuestionID,
) (*model.SessionAnswer, error) {
	query := `
		SELECT id, session_id, question_id, answer_text, answered_at
		FROM session_answers
		WHERE session_id = ? AND question_id = ?
	`

	var answer model.SessionAnswer
	var sessionIDStr string
	var questionIDInt int

	err := r.db.Conn.QueryRowContext(ctx, query, sessionID.String(), questionID.Int()).Scan(
		&answer.ID,
		&sessionIDStr,
		&questionIDInt,
		&answer.AnswerText,
		&answer.AnsweredAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewDomainError(
			errors.ANSWER_NOT_FOUND,
			"answer not found",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find session answer: %w", err)
	}

	// SessionIDとQuestionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	answer.SessionID = sid

	qid, err := valueobject.NewQuestionID(questionIDInt)
	if err != nil {
		return nil, err
	}
	answer.QuestionID = qid

	return &answer, nil
}

// GetRandomAnswers は、ランダムな過去回答を取得する（クイズのダミー選択肢用）
func (r *MySQLSessionAnswerRepository) GetRandomAnswers(
	ctx context.Context,
	questionID valueobject.QuestionID,
	limit int,
) ([]string, error) {
	query := `
		SELECT DISTINCT answer_text
		FROM session_answers
		WHERE question_id = ?
		ORDER BY RAND()
		LIMIT ?
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, questionID.Int(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get random answers: %w", err)
	}
	defer rows.Close()

	var answers []string
	for rows.Next() {
		var answer string
		if err := rows.Scan(&answer); err != nil {
			return nil, fmt.Errorf("failed to scan answer: %w", err)
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
