package service

import (
	"context"
	"math/rand"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// QuizService は、クイズ生成に関するドメインサービス
type QuizService struct {
	sessionAnswerRepo repository.SessionAnswerRepository
	quizRepo          repository.QuizQuestionRepository
}

// NewQuizService は、新しいQuizServiceを作成する
func NewQuizService(
	sessionAnswerRepo repository.SessionAnswerRepository,
	quizRepo repository.QuizQuestionRepository,
) *QuizService {
	return &QuizService{
		sessionAnswerRepo: sessionAnswerRepo,
		quizRepo:          quizRepo,
	}
}

// GenerateQuiz は、セッションの回答からクイズを生成する
func (s *QuizService) GenerateQuiz(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
) (*model.QuizQuestion, error) {
	// すでにクイズが存在する場合は、それを返す
	existingQuiz, err := s.quizRepo.FindBySessionIDAndPlaceID(ctx, sessionID, placeID)
	if err == nil && existingQuiz != nil {
		return existingQuiz, nil
	}

	// セッションの全回答を取得
	answers, err := s.sessionAnswerRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, errors.NewDomainError(
			errors.ANSWER_NOT_FOUND,
			"failed to get session answers",
			err,
		)
	}

	if len(answers) == 0 {
		return nil, errors.NewDomainError(
			errors.ANSWER_REQUIRED,
			"no answers found for session",
			nil,
		)
	}

	// ランダムに質問を選ぶ
	selectedAnswer := answers[rand.Intn(len(answers))]
	correctAnswer := selectedAnswer.AnswerText

	// ダミー選択肢を生成
	dummyOptions := s.generateDummyOptions(ctx, sessionID, answers, selectedAnswer.QuestionID, correctAnswer)

	// 4つの選択肢を作成（正解 + ダミー3つ）
	allOptions := append([]string{correctAnswer}, dummyOptions...)

	// 選択肢をシャッフル
	rand.Shuffle(len(allOptions), func(i, j int) {
		allOptions[i], allOptions[j] = allOptions[j], allOptions[i]
	})

	// 正解のインデックスを見つける
	var answerIndex int
	for i, opt := range allOptions {
		if opt == correctAnswer {
			answerIndex = i
			break
		}
	}

	// 配列にコピー
	var options [4]string
	copy(options[:], allOptions)

	// 質問文を生成
	questionText := s.generateQuestionText(selectedAnswer.QuestionID)

	// クイズを作成
	quiz := model.NewQuizQuestion(
		sessionID,
		placeID,
		questionText,
		options,
		answerIndex,
	)

	// クイズを保存
	if err := s.quizRepo.Save(ctx, quiz); err != nil {
		return nil, errors.NewDomainError(
			errors.INVALID_INPUT,
			"failed to save quiz",
			err,
		)
	}

	return quiz, nil
}

// generateDummyOptions は、ダミー選択肢を3つ生成する
func (s *QuizService) generateDummyOptions(
	ctx context.Context,
	sessionID valueobject.SessionID,
	allAnswers []*model.SessionAnswer,
	questionID valueobject.QuestionID,
	correctAnswer string,
) []string {
	dummies := make([]string, 0, 3)

	// 1. プレイヤーの別回答から取得
	for _, ans := range allAnswers {
		if !ans.QuestionID.Equals(questionID) && ans.AnswerText != correctAnswer {
			dummies = append(dummies, ans.AnswerText)
			if len(dummies) == 3 {
				return dummies
			}
		}
	}

	// 2. 過去プレイヤーの回答から取得
	randomAnswers, err := s.sessionAnswerRepo.GetRandomAnswers(ctx, questionID, 10)
	if err == nil {
		for _, ans := range randomAnswers {
			if ans != correctAnswer && !contains(dummies, ans) {
				dummies = append(dummies, ans)
				if len(dummies) == 3 {
					return dummies
				}
			}
		}
	}

	// 3. システム汎用回答で埋める
	genericAnswers := getGenericAnswers(questionID)
	for _, ans := range genericAnswers {
		if ans != correctAnswer && !contains(dummies, ans) {
			dummies = append(dummies, ans)
			if len(dummies) == 3 {
				break
			}
		}
	}

	return dummies
}

// generateQuestionText は、質問IDから質問文を生成する
func (s *QuizService) generateQuestionText(questionID valueobject.QuestionID) string {
	questions := map[int]string{
		1:  "あなたが小学生の頃、夢中になっていたことは？",
		2:  "あなたが中学生の頃、尊敬していた人は？",
		3:  "あなたが時間を忘れて没頭したことは？",
		4:  "あなたが手放したい能力は？",
		5:  "あなたが10億円と引き換えにしたくない能力は？",
		6:  "あなたが人生最後の日に話したい人と話題は？",
		7:  "あなたがお金や時間を気にせずやりたいことは？",
		8:  "あなたが住みたい場所は？",
		9:  "あなたが人生の最期に達成したいことは？",
		10: "あなたの名前は？",
	}
	return questions[questionID.Int()]
}

// getGenericAnswers は、質問IDに対応する汎用回答を返す
func getGenericAnswers(questionID valueobject.QuestionID) []string {
	genericMap := map[int][]string{
		1:  {"読書", "スポーツ", "ゲーム", "絵を描くこと"},
		2:  {"先生", "親", "友人", "スポーツ選手"},
		3:  {"音楽", "プログラミング", "料理", "旅行"},
		4:  {"完璧主義", "心配性", "優柔不断", "人見知り"},
		5:  {"家族との時間", "健康", "自由", "友人"},
		6:  {"親", "友人", "恩師", "パートナー"},
		7:  {"世界旅行", "起業", "芸術活動", "社会貢献"},
		8:  {"海辺", "山", "都会", "田舎"},
		9:  {"家族の幸せ", "夢の実現", "社会貢献", "自己実現"},
		10: {"田中", "佐藤", "鈴木", "高橋"},
	}
	return genericMap[questionID.Int()]
}

// contains は、スライスに要素が含まれているかチェックする
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
