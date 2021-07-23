package http

import (
	"encoding/json"
	"net/http"

	"gregvader/triple-tac-toe/database/database_type"
	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/dto"

	"gorm.io/gorm"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandler(usecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{AuthUsecase: usecase}
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registrationDTO dto.RegisterDTO
	trxHandle := r.Context().Value(database_type.Transactional).(*gorm.DB)

	if parseErr := json.NewDecoder(r.Body).Decode(&registrationDTO); parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	authInfo, regErr := handler.AuthUsecase.WithTrx(trxHandle).Register(registrationDTO)
	if regErr != nil {
		trxHandle.Rollback()
		http.Error(w, regErr.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authInfo)
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	trxHandle := r.Context().Value(database_type.Transactional).(*gorm.DB)
	var loginDTO dto.LoginDTO

	if parseErr := json.NewDecoder(r.Body).Decode(&loginDTO); parseErr != nil {
		http.Error(w, parseErr.Error(), http.StatusBadRequest)
		return
	}

	authInfo, authErr := handler.AuthUsecase.WithTrx(trxHandle).Login(loginDTO)
	if authErr != nil {
		http.Error(w, authErr.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authInfo)
}
