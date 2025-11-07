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

// MySQLQuizQuestionRepository は、MySQLを使用したQuizQuestionRepositoryの実装
type MySQLQuizQuestionRepository struct {
	db *database.DB
}

// NewMySQLQuizQuestionRepository は、新しいMySQLQuizQuestionRepositoryを作成する
func NewMySQLQuizQuestionRepository(db *database.DB) *MySQLQuizQuestionRepository {
	return &MySQLQuizQuestionRepository{db: db}
}

// Save は、クイズ問題を保存する
func (r *MySQLQuizQuestionRepository) Save(ctx context.Context, quiz *model.QuizQuestion) error {
	query := `
		INSERT INTO quiz_questions (quiz_id, session_id, place_id, question_text, option_1, option_2, option_3, option_4, answer_index, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Conn.ExecContext(
		ctx,
		query,
		quiz.QuizID.String(),
		quiz.SessionID.String(),
		quiz.PlaceID.String(),
		quiz.QuestionText,
		quiz.Options[0],
		quiz.Options[1],
		quiz.Options[2],
		quiz.Options[3],
		quiz.AnswerIndex,
		quiz.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save quiz question: %w", err)
	}

	return nil
}

// FindByID は、クイズIDでクイズ問題を取得する
func (r *MySQLQuizQuestionRepository) FindByID(ctx context.Context, quizID valueobject.QuizID) (*model.QuizQuestion, error) {
	query := `
		SELECT quiz_id, session_id, place_id, question_text, option_1, option_2, option_3, option_4, answer_index, created_at
		FROM quiz_questions
		WHERE quiz_id = ?
	`

	quiz, err := r.scanQuizQuestionRow(r.db.Conn.QueryRowContext(ctx, query, quizID.String()))

	if err == sql.ErrNoRows {
		return nil, errors.NewDomainError(
			errors.QUIZ_NOT_FOUND,
			"quiz question not found",
			err,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find quiz question: %w", err)
	}

	return quiz, nil
}

// FindBySessionIDAndPlaceID は、セッションIDと場所IDでクイズ問題を取得する
func (r *MySQLQuizQuestionRepository) FindBySessionIDAndPlaceID(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
) (*model.QuizQuestion, error) {
	query := `
		SELECT quiz_id, session_id, place_id, question_text, option_1, option_2, option_3, option_4, answer_index, created_at
		FROM quiz_questions
		WHERE session_id = ? AND place_id = ?
	`

	quiz, err := r.scanQuizQuestionRow(r.db.Conn.QueryRowContext(ctx, query, sessionID.String(), placeID.String()))

	if err == sql.ErrNoRows {
		return nil, nil // クイズがまだ生成されていない場合はnilを返す
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find quiz question: %w", err)
	}

	return quiz, nil
}

// FindBySessionID は、セッションIDで全クイズ問題を取得する
func (r *MySQLQuizQuestionRepository) FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) ([]*model.QuizQuestion, error) {
	query := `
		SELECT quiz_id, session_id, place_id, question_text, option_1, option_2, option_3, option_4, answer_index, created_at
		FROM quiz_questions
		WHERE session_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.Conn.QueryContext(ctx, query, sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find quiz questions: %w", err)
	}
	defer rows.Close()

	var quizzes []*model.QuizQuestion
	for rows.Next() {
		quiz, err := r.scanQuizQuestion(rows)
		if err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

// scanQuizQuestion は、rowsからQuizQuestionをスキャンする
func (r *MySQLQuizQuestionRepository) scanQuizQuestion(rows *sql.Rows) (*model.QuizQuestion, error) {
	var quiz model.QuizQuestion
	var quizIDStr string
	var sessionIDStr string
	var placeIDStr string
	var opt1, opt2, opt3, opt4 string

	err := rows.Scan(
		&quizIDStr,
		&sessionIDStr,
		&placeIDStr,
		&quiz.QuestionText,
		&opt1,
		&opt2,
		&opt3,
		&opt4,
		&quiz.AnswerIndex,
		&quiz.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan quiz question: %w", err)
	}

	// QuizIDを復元
	qid, err := valueobject.NewQuizIDFromString(quizIDStr)
	if err != nil {
		return nil, err
	}
	quiz.QuizID = qid

	// SessionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	quiz.SessionID = sid

	// PlaceIDを復元
	pid, err := valueobject.NewPlaceID(placeIDStr)
	if err != nil {
		return nil, err
	}
	quiz.PlaceID = pid

	// Optionsを設定
	quiz.Options = [4]string{opt1, opt2, opt3, opt4}

	return &quiz, nil
}

// scanQuizQuestionRow は、rowからQuizQuestionをスキャンする
func (r *MySQLQuizQuestionRepository) scanQuizQuestionRow(row *sql.Row) (*model.QuizQuestion, error) {
	var quiz model.QuizQuestion
	var quizIDStr string
	var sessionIDStr string
	var placeIDStr string
	var opt1, opt2, opt3, opt4 string

	err := row.Scan(
		&quizIDStr,
		&sessionIDStr,
		&placeIDStr,
		&quiz.QuestionText,
		&opt1,
		&opt2,
		&opt3,
		&opt4,
		&quiz.AnswerIndex,
		&quiz.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// QuizIDを復元
	qid, err := valueobject.NewQuizIDFromString(quizIDStr)
	if err != nil {
		return nil, err
	}
	quiz.QuizID = qid

	// SessionIDを復元
	sid, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		return nil, err
	}
	quiz.SessionID = sid

	// PlaceIDを復元
	pid, err := valueobject.NewPlaceID(placeIDStr)
	if err != nil {
		return nil, err
	}
	quiz.PlaceID = pid

	// Optionsを設定
	quiz.Options = [4]string{opt1, opt2, opt3, opt4}

	return &quiz, nil
}
