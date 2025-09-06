package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

type VAPIDHandler struct {
	vapidUseCase *usecase.VAPIDUseCase
}

func NewVAPIDHandler(vapidUseCase *usecase.VAPIDUseCase) *VAPIDHandler {
	return &VAPIDHandler{
		vapidUseCase: vapidUseCase,
	}
}

func (vh *VAPIDHandler) GetPublicKey(w http.ResponseWriter, r *http.Request) {
	result := vh.vapidUseCase.GetPublicKey()

	response := dto.VAPIDPublicKeyResponse{
		PublicKey: result.PublicKey,
		Success:   result.Success,
		Message:   result.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(response)
}