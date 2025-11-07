package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
)

// VAPID鍵ペアを生成するCLIツール
// Usage: go run cmd/generate-vapid/main.go
func main() {
	// NIST P-256 曲線を使用して秘密鍵を生成
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}

	// 秘密鍵をBase64 URL-safeエンコード
	privateKeyBytes := privateKey.D.Bytes()
	// 32バイトにパディング
	privateKeyBytes32 := make([]byte, 32)
	copy(privateKeyBytes32[32-len(privateKeyBytes):], privateKeyBytes)
	privateKeyB64 := base64.RawURLEncoding.EncodeToString(privateKeyBytes32)

	// 公開鍵をuncompressed形式で取得（04 || X || Y）
	publicKeyBytes := elliptic.Marshal(elliptic.P256(), privateKey.PublicKey.X, privateKey.PublicKey.Y)
	publicKeyB64 := base64.RawURLEncoding.EncodeToString(publicKeyBytes)

	// 結果を出力
	fmt.Println("========== VAPID Key Pair Generated ==========")
	fmt.Println()
	fmt.Println("Public Key (Base64 URL-safe):")
	fmt.Println(publicKeyB64)
	fmt.Println()
	fmt.Println("Private Key (Base64 URL-safe):")
	fmt.Println(privateKeyB64)
	fmt.Println()
	fmt.Println("===============================================")
	fmt.Println()
	fmt.Println("Add these to your .env file:")
	fmt.Println()
	fmt.Printf("VAPID_PUBLIC_KEY=%s\n", publicKeyB64)
	fmt.Printf("VAPID_PRIVATE_KEY=%s\n", privateKeyB64)
	fmt.Println()
	fmt.Println("Or export them as environment variables:")
	fmt.Println()
	fmt.Printf("export VAPID_PUBLIC_KEY=\"%s\"\n", publicKeyB64)
	fmt.Printf("export VAPID_PRIVATE_KEY=\"%s\"\n", privateKeyB64)
	fmt.Println()
}

// validateKeys は、生成された鍵ペアを検証する（オプション）
func validateKeys(publicKeyB64, privateKeyB64 string) error {
	// 秘密鍵をデコード
	privateKeyBytes, err := base64.RawURLEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	// 秘密鍵をECDSA形式に復元
	d := new(big.Int).SetBytes(privateKeyBytes)
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(d.Bytes())

	// 公開鍵をデコード
	publicKeyBytes, err := base64.RawURLEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// 公開鍵を復元
	expectedX, expectedY := elliptic.Unmarshal(curve, publicKeyBytes)
	if expectedX == nil || expectedY == nil {
		return fmt.Errorf("failed to unmarshal public key")
	}

	// 一致確認
	if x.Cmp(expectedX) != 0 || y.Cmp(expectedY) != 0 {
		return fmt.Errorf("key pair validation failed: public key does not match private key")
	}

	return nil
}
