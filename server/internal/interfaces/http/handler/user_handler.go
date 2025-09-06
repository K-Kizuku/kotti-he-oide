package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUseCase.GetAllUsers(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := dto.ToUsersResponse(users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetUser(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := dto.ToUserResponse(user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := dto.ToUserResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.DeleteUser(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	domainErr, ok := err.(*errors.DomainError)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var statusCode int
	switch domainErr.Code {
	case errors.ErrUserNotFound.Code:
		statusCode = http.StatusNotFound
	case errors.ErrEmailAlreadyExist.Code:
		statusCode = http.StatusConflict
	case errors.ErrInvalidUserID.Code, errors.ErrInvalidEmail.Code:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": domainErr.Message,
		"code":  domainErr.Code,
	})
}
