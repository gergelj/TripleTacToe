package http

import (
	"gregvader/triple-tac-toe/domain"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(usecase domain.UserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: usecase}
}
