package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

    pb "github.com/K-Kizuku/kotti-he-oide/internal/gen/image_recognition/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MLHandler struct {
	grpcAddr string
}

func NewMLHandler() *MLHandler {
	addr := os.Getenv("IMAGE_RECOGNITION_GRPC_ADDR")
	if addr == "" {
		addr = "127.0.0.1:50051"
	}
	return &MLHandler{grpcAddr: addr}
}

// GET /api/ml/hello?name=world
func (h *MLHandler) HelloProxy(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}

	// gRPC ダイアル（開発用に Insecure）
	conn, err := grpc.Dial(h.grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "failed to dial gRPC backend", http.StatusBadGateway)
		return
	}
	defer conn.Close()

	client := pb.NewImageRecognitionServiceClient(conn)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := client.Hello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		http.Error(w, "gRPC call failed", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": resp.GetMessage(),
		"backend": h.grpcAddr,
	})
}

// POST /api/ml/recognize
// - multipart/form-data: フィールド名は `image` または `file`
// - クエリまたはフォームで `threshold` を任意指定（0.0-1.0）
func (h *MLHandler) RecognizeImageProxy(w http.ResponseWriter, r *http.Request) {
	// 入力の Content-Type に応じて画像バイトを取得
	var imgBytes []byte

	// まずは multipart/form-data を試す
	ct := r.Header.Get("Content-Type")
	if ct != "" && (ct == "multipart/form-data" || len(ct) >= len("multipart/form-data") && ct[:len("multipart/form-data")] == "multipart/form-data") {
		// 10MB まで一時メモリに展開（超える場合は一時ファイル）
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
			http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			// 代替フィールド名 `file`
			file, _, err = r.FormFile("file")
		}
		if err == nil && file != nil {
			defer file.Close()
			buf, readErr := io.ReadAll(file)
			if readErr != nil {
				http.Error(w, "failed to read uploaded file", http.StatusBadRequest)
				return
			}
			imgBytes = buf
		}
	}

	// multipart で取れなかった場合、リクエストボディをそのまま画像として扱う
	if len(imgBytes) == 0 {
		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil || len(bodyBytes) == 0 {
			http.Error(w, "image data is required", http.StatusBadRequest)
			return
		}
		imgBytes = bodyBytes
	}

	// しきい値の取得（クエリ優先、なければフォーム値）
	var thresholdPtr *float32
	if tStr := r.URL.Query().Get("threshold"); tStr != "" {
		if tVal, err := strconv.ParseFloat(tStr, 32); err == nil {
			tv := float32(tVal)
			thresholdPtr = &tv
		}
	} else if r.PostForm != nil {
		if tStr := r.PostForm.Get("threshold"); tStr != "" {
			if tVal, err := strconv.ParseFloat(tStr, 32); err == nil {
				tv := float32(tVal)
				thresholdPtr = &tv
			}
		}
	}

	// gRPC ダイアル（開発用に Insecure）
	conn, err := grpc.Dial(h.grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "failed to dial gRPC backend", http.StatusBadGateway)
		return
	}
	defer conn.Close()

	client := pb.NewImageRecognitionServiceClient(conn)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req := &pb.RecognizeImageRequest{ImageData: imgBytes}
	if thresholdPtr != nil {
		// proto3 optional 対応（*float32 をセット）
		req.Threshold = thresholdPtr
	}

	resp, err := client.RecognizeImage(ctx, req)
	if err != nil {
		http.Error(w, "gRPC call failed", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"is_match":         resp.GetIsMatch(),
		"similarity_score": resp.GetSimilarityScore(),
		"error_message":    resp.GetErrorMessage(),
		"backend":          h.grpcAddr,
	})
}
