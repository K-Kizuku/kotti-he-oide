package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
	"github.com/SherClockHolmes/webpush-go"
)

// PushNotificationUseCase は、プッシュ通知に関するユースケース
type PushNotificationUseCase interface {
	// Subscribe は、プッシュ通知のサブスクリプションを登録する
	Subscribe(ctx context.Context, sessionID valueobject.SessionID, endpoint string, p256dh string, auth string) (*model.PushSubscription, error)

	// Unsubscribe は、プッシュ通知のサブスクリプションを削除する
	Unsubscribe(ctx context.Context, subscriptionID valueobject.SubscriptionID) error

	// SendPushNotification は、特定のセッションに即時プッシュ通知を送信する
	SendPushNotification(ctx context.Context, sessionID valueobject.SessionID, title string, message string) error

	// GetVAPIDPublicKey は、VAPID公開鍵を取得する
	GetVAPIDPublicKey(ctx context.Context) (string, error)
}

// pushNotificationUseCase は、PushNotificationUseCaseの実装
type pushNotificationUseCase struct {
	subscriptionRepo repository.PushSubscriptionRepository
	logRepo          repository.PushLogRepository
	vapidService     service.VAPIDService
}

// NewPushNotificationUseCase は、新しいPushNotificationUseCaseを作成する
func NewPushNotificationUseCase(
	subscriptionRepo repository.PushSubscriptionRepository,
	logRepo repository.PushLogRepository,
	vapidService service.VAPIDService,
) PushNotificationUseCase {
	return &pushNotificationUseCase{
		subscriptionRepo: subscriptionRepo,
		logRepo:          logRepo,
		vapidService:     vapidService,
	}
}

// Subscribe は、プッシュ通知のサブスクリプションを登録する
func (u *pushNotificationUseCase) Subscribe(ctx context.Context, sessionID valueobject.SessionID, endpoint string, p256dh string, auth string) (*model.PushSubscription, error) {
	// Value Objectsの作成
	pushEndpoint, err := valueobject.NewPushEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	keys, err := valueobject.NewPushKeys(p256dh, auth)
	if err != nil {
		return nil, err
	}

	// 既存のアクティブなサブスクリプションがあれば無効化
	existingSubscription, err := u.subscriptionRepo.FindBySessionID(ctx, sessionID)
	if err == nil && existingSubscription != nil {
		_ = existingSubscription.Deactivate()
		_ = u.subscriptionRepo.Update(ctx, existingSubscription)
	}

	// 新しいサブスクリプションを作成
	subscription := model.NewPushSubscription(sessionID, pushEndpoint, keys)

	// 保存
	if err := u.subscriptionRepo.Save(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

// Unsubscribe は、プッシュ通知のサブスクリプションを削除する
func (u *pushNotificationUseCase) Unsubscribe(ctx context.Context, subscriptionID valueobject.SubscriptionID) error {
	// サブスクリプションの存在確認
	subscription, err := u.subscriptionRepo.FindByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	// 無効化
	if err := subscription.Deactivate(); err != nil {
		return err
	}

	// 更新（論理削除）
	if err := u.subscriptionRepo.Update(ctx, subscription); err != nil {
		return err
	}

	return nil
}

// SendPushNotification は、特定のセッションに即時プッシュ通知を送信する
func (u *pushNotificationUseCase) SendPushNotification(ctx context.Context, sessionID valueobject.SessionID, title string, message string) error {
	// サブスクリプションを取得
	subscription, err := u.subscriptionRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return errors.NewDomainError(
			errors.NOT_FOUND,
			"no active push subscription found for this session",
			err,
		)
	}

	// サブスクリプションがアクティブか確認
	if !subscription.CanReceivePush() {
		return errors.NewDomainError(
			errors.INVALID_STATE,
			"push subscription is not active",
			nil,
		)
	}

	// プッシュ通知のペイロードを作成
	payload := map[string]interface{}{
		"title":   title,
		"message": message,
		"icon":    "/icon.png", // アイコンパス（フロントエンドのpublic/icon.png）
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal push payload: %w", err)
	}

	// Web Push APIで送信
	webpushSubscription := &webpush.Subscription{
		Endpoint: subscription.Endpoint.String(),
		Keys: webpush.Keys{
			P256dh: subscription.Keys.P256dh(),
			Auth:   subscription.Keys.Auth(),
		},
	}

	options := &webpush.Options{
		Subscriber:      "mailto:noreply@kotti-he-oide.com", // 管理者メールアドレス
		VAPIDPublicKey:  u.vapidService.GetPublicKey(),
		VAPIDPrivateKey: u.vapidService.GetPrivateKey(),
		TTL:             30, // 30秒のTTL（ホラー演出用なので短め）
	}

	resp, err := webpush.SendNotification(payloadBytes, webpushSubscription, options)
	if resp != nil {
		defer resp.Body.Close()
	}

	// 送信結果をログに記録
	if err != nil {
		// エラーの場合
		statusCode := 0
		if resp != nil {
			statusCode = resp.StatusCode
		}

		pushLog := model.NewPushLogWithError(
			subscription.SubscriptionID,
			sessionID,
			title,
			message,
			statusCode,
			err.Error(),
		)

		_ = u.logRepo.Save(ctx, pushLog)

		// 404/410エラーの場合はサブスクリプションを無効化
		if statusCode == 404 || statusCode == 410 {
			_ = subscription.Deactivate()
			_ = u.subscriptionRepo.Update(ctx, subscription)
		}

		return fmt.Errorf("failed to send push notification: %w", err)
	}

	// 成功の場合
	pushLog := model.NewPushLog(
		subscription.SubscriptionID,
		sessionID,
		title,
		message,
		resp.StatusCode,
	)

	if err := u.logRepo.Save(ctx, pushLog); err != nil {
		// ログ保存失敗は致命的ではないのでエラーを記録するだけ
		fmt.Printf("failed to save push log: %v\n", err)
	}

	return nil
}

// GetVAPIDPublicKey は、VAPID公開鍵を取得する
func (u *pushNotificationUseCase) GetVAPIDPublicKey(ctx context.Context) (string, error) {
	return u.vapidService.GetPublicKey(), nil
}
