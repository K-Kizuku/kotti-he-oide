-- プッシュ通知サブスクリプションテーブル
CREATE TABLE IF NOT EXISTS push_subscriptions (
    subscription_id CHAR(36) PRIMARY KEY,
    session_id CHAR(36) NOT NULL,
    endpoint TEXT NOT NULL,
    p256dh_key VARCHAR(255) NOT NULL,
    auth_key VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_push_subscriptions_session_id (session_id),
    INDEX idx_push_subscriptions_is_active (is_active),
    FOREIGN KEY (session_id) REFERENCES sessions(session_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- プッシュ通知送信ログテーブル
CREATE TABLE IF NOT EXISTS push_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    subscription_id CHAR(36) NOT NULL,
    session_id CHAR(36) NOT NULL,
    title VARCHAR(255),
    message TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    status_code INT,
    error_message TEXT,
    sent_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_push_logs_subscription_id (subscription_id),
    INDEX idx_push_logs_session_id (session_id),
    INDEX idx_push_logs_sent_at (sent_at),
    FOREIGN KEY (subscription_id) REFERENCES push_subscriptions(subscription_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
