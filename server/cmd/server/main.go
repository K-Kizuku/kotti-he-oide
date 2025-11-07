package main

import (
	"context"
	"log"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/database"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/persistence"
	"github.com/K-Kizuku/kotti-he-oide/internal/infrastructure/voicevox"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/handler"
	"github.com/K-Kizuku/kotti-he-oide/pkg/config"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// データベース接続
	ctx := context.Background()
	db, err := database.NewMySQLDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// リポジトリの初期化
	sessionRepo := persistence.NewMySQLSessionRepository(db)
	sessionAnswerRepo := persistence.NewMySQLSessionAnswerRepository(db)
	s6ProgressRepo := persistence.NewMySQLS6ProgressRepository(db)
	quizRepo := persistence.NewMySQLQuizQuestionRepository(db)
	messageRepo := persistence.NewMySQLPlayerMessageRepository(db)
	pushSubscriptionRepo := persistence.NewMySQLPushSubscriptionRepository(db)
	pushLogRepo := persistence.NewMySQLPushLogRepository(db)

	// ドメインサービスの初期化
	sessionService := service.NewSessionService(sessionRepo)
	quizService := service.NewQuizService(sessionAnswerRepo, quizRepo)
	s6Service := service.NewS6Service(s6ProgressRepo, sessionRepo)

	// VAPID サービスの初期化
	vapidService, err := service.NewVAPIDService()
	if err != nil {
		log.Fatal("Failed to initialize VAPID service:", err)
	}

	// VOICEVOXクライアントの初期化
	voicevoxClient := voicevox.NewClient(cfg.VoicevoxAPIURL)

	// ユースケースの初期化
	sessionUseCase := usecase.NewSessionUseCase(sessionRepo, sessionService, cfg.SessionTTLMinutes)
	s4AnswerUseCase := usecase.NewS4AnswerUseCase(sessionAnswerRepo, sessionService)
	s6UseCase := usecase.NewS6UseCase(s6ProgressRepo, quizService, s6Service, sessionService)
	messageUseCase := usecase.NewMessageUseCase(messageRepo, sessionService)
	pushNotificationUseCase := usecase.NewPushNotificationUseCase(pushSubscriptionRepo, pushLogRepo, vapidService)
	voiceUseCase := usecase.NewVoiceUseCase(voicevoxClient, cfg)

	// ハンドラーの初期化
	healthHandler := handler.NewHealthHandler()
	sessionHandler := handler.NewSessionHandler(sessionUseCase)
	answerHandler := handler.NewAnswerHandler(s4AnswerUseCase)
	s6Handler := handler.NewS6Handler(s6UseCase)
	messageHandler := handler.NewMessageHandler(messageUseCase)
	pushNotificationHandler := handler.NewPushNotificationHandler(pushNotificationUseCase)
	voiceHandler := handler.NewVoiceHandler(voiceUseCase)

	// ルーティング設定
	mux := http.NewServeMux()

	// CORS対応（開発用）
	corsHandler := enableCORS(mux)

	// Health check
	mux.HandleFunc("GET /api/healthz", healthHandler.HealthCheck)

	// Session API
	mux.HandleFunc("POST /api/session", sessionHandler.CreateSession)
	mux.HandleFunc("GET /api/session/{session_id}", sessionHandler.GetSession)
	mux.HandleFunc("POST /api/session/{session_id}/s6/start", sessionHandler.StartS6)

	// S4 Answer API
	mux.HandleFunc("POST /api/session/{session_id}/answers", answerHandler.SaveAnswer)
	mux.HandleFunc("GET /api/session/{session_id}/answers", answerHandler.GetAnswers)

	// S6 Progress API
	mux.HandleFunc("POST /api/session/{session_id}/s6/initialize", s6Handler.InitializeProgress)
	mux.HandleFunc("POST /api/session/{session_id}/s6/verify-location", s6Handler.VerifyLocation)
	mux.HandleFunc("GET /api/session/{session_id}/s6/quiz/{place_id}", s6Handler.GetQuiz)
	mux.HandleFunc("POST /api/session/{session_id}/s6/answer", s6Handler.SubmitAnswer)
	mux.HandleFunc("GET /api/session/{session_id}/s6/progress", s6Handler.GetProgress)

	// Message API
	mux.HandleFunc("POST /api/session/{session_id}/message", messageHandler.SaveMessage)
	mux.HandleFunc("GET /api/messages", messageHandler.GetMessages)

	// Push Notification API
	mux.HandleFunc("GET /api/push/vapid-public-key", pushNotificationHandler.GetVAPIDPublicKey)
	mux.HandleFunc("POST /api/push/subscribe", pushNotificationHandler.Subscribe)
	mux.HandleFunc("DELETE /api/push/subscriptions/{subscription_id}", pushNotificationHandler.Unsubscribe)
	mux.HandleFunc("POST /api/push/send/{session_id}", pushNotificationHandler.SendPush)

	// Voice API
	mux.HandleFunc("POST /api/voice/generate", voiceHandler.GenerateVoice)

	// 静的ファイル配信（音声ファイル）
	audioFileServer := http.FileServer(http.Dir(cfg.AudioOutputDir))
	mux.Handle("/audio/", http.StripPrefix("/audio/", audioFileServer))

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("VOICEVOX API URL: %s", cfg.VoicevoxAPIURL)
	log.Printf("Audio output directory: %s", cfg.AudioOutputDir)
	if err := http.ListenAndServe(":"+cfg.ServerPort, corsHandler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// enableCORS は、CORSを有効にするミドルウェア
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
