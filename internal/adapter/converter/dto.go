package converter

import (
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/dto"
	"time"
)

func UserRegisterToModel(userDto *dto.UserRegister, userRole entity.UserRole, pass string) (*entity.User, error) {

	compType := entity.CompanyTypes[userDto.CompanyType]
	return &entity.User{
		UserRole:       userRole,
		Name:           userDto.Name,
		Surname:        userDto.Surname,
		Email:          userDto.Email,
		PhoneNumber:    userDto.PhoneNumber,
		CompanyName:    userDto.CompanyName,
		CompanyType:    compType,
		CompanyWebSite: userDto.CompanyWebSite,
		PasswordHash:   pass,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
