package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "time"

    pb "github.com/K-Kizuku/kotti-he-oide/server/internal/gen/image_recognition/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type MLHandler struct{
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

