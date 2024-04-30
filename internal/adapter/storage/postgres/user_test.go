package psql

import (
	"context"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		setup  func(db db.EngineMaker) error
		input  entity.User
		errors bool
	}{
		{
			name: "happy path",
			input: entity.User{
				Name:           "Bulut",
				Surname:        "Gocer",
				UserRole:       "customer",
				Email:          "bulut@gmail.com",
				PhoneNumber:    "1233212",
				CompanyName:    "Yuka",
				CompanyType:    "Type A",
				CompanyWebSite: "yuka.com",
				PasswordHash:   "1234321",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			errors: false,
		},
		{
			name: "empty name",
			input: entity.User{
				Name:           "",
				Surname:        "Gocer",
				UserRole:       "customer",
				Email:          "bulutcan@gmail.com",
				PhoneNumber:    "1233212",
				CompanyName:    "Yuka",
				CompanyType:    "Type A",
				CompanyWebSite: "yuka.com",
				PasswordHash:   "1234321",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			errors: true,
		},
		{
			name: "same email/phone number collision",
			input: entity.User{
				Name:           "Yanki Can",
				Surname:        "Gocer",
				UserRole:       "customer",
				Email:          "yanki@gmail.com",
				PhoneNumber:    "12332131",
				CompanyName:    "Yuka",
				CompanyType:    "Type A",
				CompanyWebSite: "yuka.com",
				PasswordHash:   "1234321",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			setup: func(db db.EngineMaker) error {
				userRepo := NewUserRepository(db)
				_, err := userRepo.Insert(context.Background(), &entity.User{
					Name:           "Yanki",
					Surname:        "Gocer",
					UserRole:       "customer",
					Email:          "yanki@gmail.com",
					PhoneNumber:    "12332131",
					CompanyName:    "Yuka",
					CompanyType:    "Type A",
					CompanyWebSite: "yuka.com",
					PasswordHash:   "1234321",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				})
				return err
			},
			errors: true,
		},
	}
	fmt.Println(testCases)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			engine := getConnection()
			assert.NotNil(t, engine)
			err := getConnection().Execute(ctx, "DELETE FROM users")
			if err != nil {
				fmt.Println("313131", err)
			}
			t.Cleanup(cleanUp)
			repo := NewUserRepository(engine)
			if tc.setup != nil {
				tc.setup(engine)
			}

			userID, err := repo.Insert(ctx, &tc.input)
			time.Sleep(sleepTime)
			if tc.errors {
				fmt.Println(err)
			}
			assert.NoError(t, err)
			err = engine.Execute(
				ctx,
				"SELECT * FROM users WHERE id = $1",
				userID,
			)

			assert.NoError(t, err)
		})
	}
}
