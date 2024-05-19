package converter

import (
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/dto"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"time"
)

func UserRegisterRequestToModel(userDto *dto.UserRegisterRequest, role, pass string) (*entity.User, error) {
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
