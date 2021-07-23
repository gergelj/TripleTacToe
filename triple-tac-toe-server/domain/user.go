package domain

import (
	"gregvader/triple-tac-toe/dto"
	"gregvader/triple-tac-toe/security"

	"gorm.io/gorm"
)

type User struct {
	Id       uint64 `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null" valid:"stringlength(2|20),required"`
	Password string `gorm:"not null" valid:"stringlength(8|50),required"`
}

type UserRepository interface {
	Create(User User) (User, error)
	FindByUsername(username string) (User, error)
	WithTrx(trxHandle *gorm.DB) UserRepository
}

type UserUsecase interface {
	FindByUsername(username string) (User, error)
	WithTrx(trxHanle *gorm.DB) UserUsecase
}

type AuthUsecase interface {
	Register(registerDTO dto.RegisterDTO) (security.AuthInfo, error)
	Login(loginDTO dto.LoginDTO) (security.AuthInfo, error)
	WithTrx(trxHandle *gorm.DB) AuthUsecase
}
