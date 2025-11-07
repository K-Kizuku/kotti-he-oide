package dto

import (
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
)

// QuizResponse は、クイズ情報のレスポンスDTO
type QuizResponse struct {
	QuizID       string   `json:"quiz_id"`
	PlaceID      string   `json:"place_id"`
	QuestionText string   `json:"question_text"`
	Options      []string `json:"options"`
}

// ToQuizResponse は、QuizQuestionモデルからQuizResponseDTOに変換する
func ToQuizResponse(quiz *model.QuizQuestion) *QuizResponse {
	return &QuizResponse{
		QuizID:       quiz.QuizID.String(),
		PlaceID:      quiz.PlaceID.String(),
		QuestionText: quiz.QuestionText,
		Options:      quiz.Options[:],
	}
}

// QuizAnswerRequest は、クイズ回答のリクエストDTO
type QuizAnswerRequest struct {
	AnswerIndex int `json:"answer_index"`
}

// QuizAnswerResponse は、クイズ回答結果のレスポンスDTO
type QuizAnswerResponse struct {
	Correct bool `json:"correct"`
}

// VerifyLocationRequest は、場所到達検証のリクエストDTO（Web選択のみ）
type VerifyLocationRequest struct {
	PlaceID string `json:"place_id"`
}

// VerifyLocationResponse は、場所到達検証のレスポンスDTO
type VerifyLocationResponse struct {
	Verified bool `json:"verified"`
}

// S6ProgressResponse は、S6進捗情報のレスポンスDTO
type S6ProgressResponse struct {
	PlaceID     string `json:"place_id"`
	Verified    bool   `json:"verified"`
	Answered    bool   `json:"answered"`
	Correct     bool   `json:"correct"`
	IsCompleted bool   `json:"is_completed"`
}

// ToS6ProgressResponse は、S6Progressモデルから S6ProgressResponseDTOに変換する
func ToS6ProgressResponse(progress *model.S6Progress) *S6ProgressResponse {
	return &S6ProgressResponse{
		PlaceID:     progress.PlaceID.String(),
		Verified:    progress.Verified,
		Answered:    progress.Answered,
		Correct:     progress.Correct,
		IsCompleted: progress.IsCompleted(),
	}
}

// ToS6ProgressResponses は、S6Progressモデルのスライスから S6ProgressResponseDTOのスライスに変換する
func ToS6ProgressResponses(progressList []*model.S6Progress) []*S6ProgressResponse {
	responses := make([]*S6ProgressResponse, len(progressList))
	for i, progress := range progressList {
		responses[i] = ToS6ProgressResponse(progress)
	}
	return responses
}
