package converter

import (
	"github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
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

func UserModelToUpdateArg(userModel *entity.User) *user.UpdateUserParams {
	changedUserParams := &user.UpdateUserParams{
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	setIfNotEmpty := func(field string, value any) {
		if value != "" {
			switch field {
			case "Name":
				changedUserParams.Name = pgtype.Text{String: value.(string), Valid: true}
			case "Email":
				changedUserParams.Email = pgtype.Text{String: value.(string), Valid: true}
			case "PhoneNumber":
				changedUserParams.PhoneNumber = pgtype.Text{String: value.(string), Valid: true}
			case "CompanyName":
				changedUserParams.CompanyName = pgtype.Text{String: value.(string), Valid: true}
			case "CompanyType":
				changedUserParams.CompanyType = pgtype.Text{String: value.(string), Valid: true}
			case "CompanyWebSite":
				changedUserParams.CompanyWebsite = pgtype.Text{String: value.(string), Valid: true}
			case "PasswordHash":
				changedUserParams.PasswordHash = pgtype.Text{String: value.(string), Valid: true}
			case "Surname":
				changedUserParams.Surname = pgtype.Text{String: value.(string), Valid: true}
			}
		}
	}

	setIfNotEmpty("Name", userModel.Name)
	setIfNotEmpty("Email", userModel.Email)
	setIfNotEmpty("PhoneNumber", userModel.PhoneNumber)
	setIfNotEmpty("CompanyName", userModel.CompanyName)
	setIfNotEmpty("CompanyType", userModel.CompanyType)
	setIfNotEmpty("CompanyWebSite", userModel.CompanyWebSite)
	setIfNotEmpty("PasswordHash", userModel.PasswordHash)
	setIfNotEmpty("Surname", userModel.Surname)

	return changedUserParams
}
