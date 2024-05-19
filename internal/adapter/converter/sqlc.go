package converter

import (
	user_sql "github.com/go-matchmaker/matchmaker-server/internal/adapter/storage/postgres/sqlc/generated/user"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func ArgToUserModel(userParam *user_sql.Users) *entity.User {
	return &entity.User{
		ID:           userParam.ID,
		Role:         userParam.Role,
		Name:         userParam.Name,
		Surname:      userParam.Surname,
		Email:        userParam.Email,
		PhoneNumber:  userParam.PhoneNumber,
		PasswordHash: userParam.PasswordHash,
		CreatedAt:    userParam.CreatedAt,
		UpdatedAt:    userParam.UpdatedAt,
	}
}

func UserModelToArg(userModel *entity.User) *user_sql.Users {
	return &user_sql.Users{
		ID:           userModel.ID,
		Role:         userModel.Role,
		Name:         userModel.Name,
		Surname:      userModel.Surname,
		Email:        userModel.Email,
		PhoneNumber:  userModel.PhoneNumber,
		PasswordHash: userModel.PasswordHash,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}
}

func UserModelToUpdateArg(userModel *entity.User) *user_sql.UpdateParams {
	changedUserParams := &user_sql.UpdateParams{
		ID:        userModel.ID,
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
	setIfNotEmpty("PasswordHash", userModel.PasswordHash)
	setIfNotEmpty("Surname", userModel.Surname)

	return changedUserParams
}
