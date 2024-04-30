package psql

import (
	"context"
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
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
	//for _, tc := range testCases {
	//	t.Run(tc.name, func(t *testing.T) {
	//		ctx := context.Background()
	//
	//		err := getConnection(ctx, cfg)
	//		assert.NoError(t, err)
	//
	//		repo := repository.New(conn)
	//
	//		t.Cleanup(cleanup)
	//
	//		if tc.setup != nil {
	//			tc.setup(ctx, conn)
	//		}
	//
	//		now := time.Now().Truncate(time.Millisecond)
	//		spell, err := repo.Create(ctx, tc.input)
	//		time.Sleep(sleepTime)
	//		after := time.Now().Truncate(time.Millisecond)
	//
	//		if tc.errors {
	//			assert.Error(t, err)
	//
	//			// Ensure nothing exists in database
	//			row := conn.QueryRowContext(
	//				ctx,
	//				"SELECT id FROM spell WHERE name = $1 AND damage = $2 AND mana = $3",
	//				tc.input.Name,
	//				tc.input.Damage,
	//				tc.input.Mana,
	//			)
	//
	//			var id string
	//			err := row.Scan(&id)
	//
	//			assert.ErrorIs(t, err, sql.ErrNoRows)
	//
	//			return
	//		}
	//		// Assert no error
	//		assert.NoError(t, err)
	//
	//		// Check spell properties
	//		assert.Equal(t, spell.Name, tc.input.Name)
	//		assert.Equal(t, spell.Mana, tc.input.Mana)
	//		assert.Equal(t, spell.Damage, tc.input.Damage)
	//		assert.Equal(t, spell.CreatedAt, spell.UpdatedAt)
	//		assert.GreaterOrEqual(t, spell.CreatedAt, now)
	//		assert.LessOrEqual(t, spell.CreatedAt, after)
	//
	//		// Ensure row exists in database
	//		row := conn.QueryRowContext(
	//			ctx,
	//			"SELECT id, name, mana, damage, created_at, updated_at FROM spell WHERE id = $1",
	//			spell.ID,
	//		)
	//
	//		var rowSpell repository.Spell
	//		err = row.Scan(
	//			&rowSpell.ID,
	//			&rowSpell.Name,
	//			&rowSpell.Mana,
	//			&rowSpell.Damage,
	//			&rowSpell.CreatedAt,
	//			&rowSpell.UpdatedAt,
	//		)
	//
	//		assert.NoError(t, err)
	//		assert.Equal(t, rowSpell, spell)
	//	})
	//}

}
