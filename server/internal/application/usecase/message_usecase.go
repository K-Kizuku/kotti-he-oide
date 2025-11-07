package usecase

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// MessageUseCase は、プレイヤーメッセージ管理のユースケース
type MessageUseCase struct {
	messageRepo    repository.PlayerMessageRepository
	sessionService *service.SessionService
}

// NewMessageUseCase は、新しいMessageUseCaseを作成する
func NewMessageUseCase(
	messageRepo repository.PlayerMessageRepository,
	sessionService *service.SessionService,
) *MessageUseCase {
	return &MessageUseCase{
		messageRepo:    messageRepo,
		sessionService: sessionService,
	}
}

// SaveMessage は、プレイヤーメッセージを保存する
func (u *MessageUseCase) SaveMessage(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
	messageText string,
) error {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return err
	}

	// メッセージを作成
	message, err := model.NewPlayerMessage(sessionID, placeID, messageText)
	if err != nil {
		return err
	}

	// メッセージを保存
	return u.messageRepo.Save(ctx, message)
}

// GetMessagesByPlace は、場所IDでメッセージ一覧を取得する
func (u *MessageUseCase) GetMessagesByPlace(
	ctx context.Context,
	placeID valueobject.PlaceID,
	limit int,
) ([]*model.PlayerMessage, error) {
	return u.messageRepo.FindByPlaceID(ctx, placeID, limit)
}

// GetAllMessages は、全メッセージを取得する
func (u *MessageUseCase) GetAllMessages(ctx context.Context, limit int) ([]*model.PlayerMessage, error) {
	return u.messageRepo.FindAll(ctx, limit)
}

// GetMessageCount は、メッセージの総数を取得する
func (u *MessageUseCase) GetMessageCount(ctx context.Context) (int, error) {
	return u.messageRepo.Count(ctx)
}
