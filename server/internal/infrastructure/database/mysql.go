package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB は、MySQL接続のラッパー
type DB struct {
	Conn *sql.DB
}

// NewMySQLDB は、MySQL接続を作成する
func NewMySQLDB(ctx context.Context, databaseURL string) (*DB, error) {
	// MySQL接続を開く
	conn, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}

	// 接続プールの設定
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// 接続テスト
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to MySQL database")

	return &DB{Conn: conn}, nil
}

// Close は、データベース接続をクローズする
func (db *DB) Close() {
	db.Conn.Close()
}

// トランザクション管理用のヘルパー関数

// WithTx は、トランザクション内で関数を実行する
func (db *DB) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// トランザクション内でエラーが発生した場合はロールバック
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 関数を実行
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
