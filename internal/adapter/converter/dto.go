package converter

import (
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/dto"
	"time"
)

func UserRegisterToModel(userDto *dto.UserRegister, role, pass string) (*entity.User, error) {

	return &entity.User{
		Role:         role,
		Name:         userDto.Name,
		Surname:      userDto.Surname,
		Email:        userDto.Email,
		PhoneNumber:  userDto.PhoneNumber,
		PasswordHash: pass,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
