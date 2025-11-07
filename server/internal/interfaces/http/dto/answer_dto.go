package dto

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
)

// AnswerRequest は、回答保存のリクエストDTO
type AnswerRequest struct {
	QuestionID int    `json:"question_id"`
	AnswerText string `json:"answer_text"`
}

// AnswerResponse は、回答情報のレスポンスDTO
type AnswerResponse struct {
	QuestionID int       `json:"question_id"`
	AnswerText string    `json:"answer_text"`
	AnsweredAt time.Time `json:"answered_at"`
}

// ToAnswerResponse は、SessionAnswerモデルからAnswerResponseDTOに変換する
func ToAnswerResponse(answer *model.SessionAnswer) *AnswerResponse {
	return &AnswerResponse{
		QuestionID: answer.QuestionID.Int(),
		AnswerText: answer.AnswerText,
		AnsweredAt: answer.AnsweredAt,
	}
}

// ToAnswerResponses は、SessionAnswerモデルのスライスからAnswerResponseDTOのスライスに変換する
func ToAnswerResponses(answers []*model.SessionAnswer) []*AnswerResponse {
	responses := make([]*AnswerResponse, len(answers))
	for i, answer := range answers {
		responses[i] = ToAnswerResponse(answer)
	}
	return responses
}
