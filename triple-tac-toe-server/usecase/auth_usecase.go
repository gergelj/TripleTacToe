package usecase

import (
	"errors"

	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/dto"
	"gregvader/triple-tac-toe/mappers"
	"gregvader/triple-tac-toe/security"
	"gregvader/triple-tac-toe/security/helpers/hash"
	tokens "gregvader/triple-tac-toe/security/helpers/token"

	"gorm.io/gorm"
)

type authUsecase struct {
	userRepo domain.UserRepository
}

func (a authUsecase) WithTrx(trxHandle *gorm.DB) domain.AuthUsecase {
	a.userRepo = a.userRepo.WithTrx(trxHandle)
	return a
}

func (a authUsecase) Login(loginDTO dto.LoginDTO) (security.AuthInfo, error) {
	user, userErr := a.userRepo.FindByUsername(loginDTO.Username)
	if userErr != nil {
		return security.AuthInfo{}, errors.New("login error: user with username '" + loginDTO.Username + "' could not be found")
	}

	if !isLoginDataCorrect(user, loginDTO) {
		return security.AuthInfo{}, errors.New("login error: username and/or password is not correct")
	}

	return tokens.CreateAuthInfo(user)

}

func isLoginDataCorrect(user domain.User, loginDTO dto.LoginDTO) bool {
	return (loginDTO.Username == user.Username) && (hash.CompareHashedPwWithPlainPw(user.Password, loginDTO.Password))
}

func (a authUsecase) Register(registerUserDTO dto.RegisterDTO) (security.AuthInfo, error) {
	user := mappers.MapRegistrationDTOToUser(registerUserDTO)

	regUser, createErr := a.userRepo.Create(user)
	if createErr != nil {
		return security.AuthInfo{}, createErr
	}

	return tokens.CreateAuthInfo(regUser)
}

func NewAuthUsecase(userRepo domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}
