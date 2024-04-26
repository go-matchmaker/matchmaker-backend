package converter

import (
	"github.com/bulutcan99/company-matcher/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/bulutcan99/company-matcher/internal/core/domain"
)

func ArgToUserModel(user *user.Users) *domain.User {
	return &domain.User{
		ID:             user.ID,
		UserRole:       domain.UserRole(user.UserRole),
		Name:           user.Name,
		Surname:        user.Surname,
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		CompanyName:    user.CompanyName,
		CompanyType:    int(user.CompanyType),
		CompanyWebSite: user.CompanyWebsite,
		PasswordHash:   user.PasswordHash,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func UserModelToArg(userModel *domain.User) *user.Users {
	return &user.Users{
		ID:             userModel.ID,
		UserRole:       user.UserRole(userModel.UserRole),
		Name:           userModel.Name,
		Surname:        userModel.Surname,
		Email:          userModel.Email,
		PhoneNumber:    userModel.PhoneNumber,
		CompanyName:    userModel.CompanyName,
		CompanyType:    int32(userModel.CompanyType),
		CompanyWebsite: userModel.CompanyWebSite,
		PasswordHash:   userModel.PasswordHash,
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}
}
