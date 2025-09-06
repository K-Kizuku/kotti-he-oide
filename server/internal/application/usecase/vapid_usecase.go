package usecase

import (
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
)

type GetVAPIDPublicKeyResponse struct {
	PublicKey string
	Success   bool
	Message   string
}

type VAPIDUseCase struct {
	vapidService *service.VAPIDService
}

func NewVAPIDUseCase(vapidService *service.VAPIDService) *VAPIDUseCase {
	return &VAPIDUseCase{
		vapidService: vapidService,
	}
}

func (vu *VAPIDUseCase) GetPublicKey() *GetVAPIDPublicKeyResponse {
	publicKey := vu.vapidService.GetPublicKey()

	return &GetVAPIDPublicKeyResponse{
		PublicKey: publicKey,
		Success:   true,
		Message:   "VAPID public key retrieved successfully",
	}
}
