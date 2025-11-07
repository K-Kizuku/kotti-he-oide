-- プッシュ通知関連テーブルの削除（外部キー制約を考慮して逆順）
DROP TABLE IF EXISTS push_logs;
DROP TABLE IF EXISTS push_subscriptions;
