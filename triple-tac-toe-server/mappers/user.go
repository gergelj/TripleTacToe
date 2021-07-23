package mappers

import (
	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/dto"
)

func MapRegistrationDTOToUser(registrationDTO dto.RegisterDTO) domain.User {
	return domain.User{
		Username: registrationDTO.Username,
		Password: registrationDTO.Password,
	}
}

func MapUserToUserInfoDTO(user domain.User) dto.UserInfoDTO {
	return dto.UserInfoDTO{
		Id:       user.Id,
		Username: user.Username,
	}
}
