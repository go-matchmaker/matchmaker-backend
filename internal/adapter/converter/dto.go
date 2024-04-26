package converter

import (
	"github.com/bulutcan99/company-matcher/internal/core/domain/entity"
	"github.com/bulutcan99/company-matcher/internal/core/dto"
	"github.com/google/uuid"
	"time"
)

func UserRegisterToModel(userDto *dto.UserRegister, userRole entity.UserRole, pass string) (*entity.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:             id,
		UserRole:       userRole,
		Name:           userDto.Name,
		Surname:        userDto.Surname,
		Email:          userDto.Email,
		PhoneNumber:    userDto.PhoneNumber,
		CompanyName:    userDto.CompanyName,
		CompanyType:    userDto.CompanyType,
		CompanyWebSite: userDto.CompanyWebSite,
		PasswordHash:   pass,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
