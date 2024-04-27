package converter

import (
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
)

func ArgToUserModel(user *user.Users) *entity.User {
	return &entity.User{
		ID:             user.ID,
		UserRole:       entity.UserRole(user.UserRole),
		Name:           user.Name,
		Surname:        user.Surname,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		CompanyName:    user.CompanyName,
		CompanyType:    user.CompanyType,
		CompanyWebSite: user.CompanyWebsite,
		PasswordHash:   user.PasswordHash,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func UserModelToArg(userModel *entity.User) *user.Users {
	return &user.Users{
		ID:             userModel.ID,
		UserRole:       user.UserRole(userModel.UserRole),
		Name:           userModel.Name,
		Surname:        userModel.Surname,
		Email:          userModel.Email,
		PhoneNumber:    userModel.PhoneNumber,
		CompanyName:    userModel.CompanyName,
		CompanyType:    userModel.CompanyType,
		CompanyWebsite: userModel.CompanyWebSite,
		PasswordHash:   userModel.PasswordHash,
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}
}
