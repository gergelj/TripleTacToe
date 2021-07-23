package usecase

import (
	"gregvader/triple-tac-toe/domain"

	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func (p userUsecase) WithTrx(trxHandle *gorm.DB) domain.UserUsecase {
	p.userRepo = p.userRepo.WithTrx(trxHandle)
	return p
}

func (p userUsecase) FindByUsername(username string) (domain.User, error) {
	return p.userRepo.FindByUsername(username)
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo}
}
