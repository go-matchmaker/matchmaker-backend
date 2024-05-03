package psql

import (
	"fmt"
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/entity"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/db"
	"github.com/go-matchmaker/matchmaker-server/internal/core/port/repository"
	"github.com/go-matchmaker/matchmaker-server/internal/core/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func getUserRepo(newDB db.EngineMaker) repository.UserPort {
	userRepo := NewUserRepository(newDB)
	return userRepo
}

func createRandomUser() *entity.User {
	pass, err := util.HashPassword(util.RandomString(15))
	if err != nil {
		log.Fatal(err)
	}

	return &entity.User{
		UserRole:       "customer",
		Name:           util.RandomOwner(),
		Surname:        util.RandomOwner(),
		Email:          util.RandomEmail(),
		PhoneNumber:    util.RandomPhoneNumber(),
		CompanyName:    util.RandomOwner(),
		CompanyType:    "Type A",
		CompanyWebSite: util.RandomWebSite(),
		PasswordHash:   pass,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
func TestCreate(t *testing.T) {
	t.Parallel()
	user := createRandomUser()

	testCases := []struct {
		name   string
		setup  func(repo repository.UserPort) error
		input  *entity.User
		errors bool
	}{
		{
			name:   "happy path",
			input:  user,
			errors: false,
		},
		{
			name:   "violates unique constraint",
			input:  user,
			errors: true,
		},
	}

	engine := getConnection()
	time.Sleep(2 * time.Second)
	require.NotNil(t, engine)
	userRepo := getUserRepo(engine)
	require.NotNil(t, userRepo)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.setup != nil {
				err := tc.setup(userRepo)
				require.NoError(t, err)
			}

			id, err := userRepo.Insert(ctx, tc.input)
			if tc.errors {
				require.NotNil(t, err)

			} else {
				require.NoError(t, err)
				require.NotNil(t, id)
			}
		})
	}
	fmt.Println("Test create user done")
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	user := createRandomUser()

	testCases := []struct {
		name   string
		setup  func(repo repository.UserPort) (*uuid.UUID, error)
		input  *entity.User
		errors bool
	}{
		{
			name: "happy path",
			setup: func(repo repository.UserPort) (*uuid.UUID, error) {
				id, err := repo.Insert(ctx, user)
				return id, err
			},
			input:  user,
			errors: false,
		},
		{
			name: "not found",
			setup: func(repo repository.UserPort) (*uuid.UUID, error) {
				id, err := uuid.NewV7()
				return &id, err
			},
			input:  user,
			errors: true,
		},
	}

	engine := getConnection()
	time.Sleep(2 * time.Second)
	require.NotNil(t, engine)
	userRepo := getUserRepo(engine)
	require.NotNil(t, userRepo)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			id, err := tc.setup(userRepo)
			require.NotNil(t, id)
			require.NoError(t, err)
			userData, err := userRepo.GetUserByID(ctx, *id)
			if tc.errors {
				require.NotNil(t, err)
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, userData)
			}
		})
	}
	fmt.Println("Test get user done")
}
func TestDeleteUser(t *testing.T) {
	t.Parallel()
	user := createRandomUser()
	testCases := []struct {
		name   string
		setup  func(repo repository.UserPort) (*entity.User, error)
		errors bool
	}{
		{
			name: "happy path",
			setup: func(repo repository.UserPort) (*entity.User, error) {
				id, err := repo.Insert(ctx, user)
				if err != nil {
					return nil, err
				}
				user.ID = *id
				return user, nil
			},
			errors: false,
		},
		{
			name: "not uuid standard",
			setup: func(repo repository.UserPort) (*entity.User, error) {
				id := "random"
				userID, err := uuid.Parse(id)
				if err != nil {
					return nil, err
				}
				nonStandardUser := &entity.User{
					ID: userID,
				}
				return nonStandardUser, nil
			},
			errors: true,
		},
	}

	engine := getConnection()
	time.Sleep(2 * time.Second)
	require.NotNil(t, engine)
	userRepo := getUserRepo(engine)
	require.NotNil(t, userRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			user, err := tc.setup(userRepo)
			if tc.errors {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				err = userRepo.DeleteUser(ctx, user.ID)

				require.NoError(t, err)

			}
		})
	}
	fmt.Println("Test delete user done")
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	user := createRandomUser()
	testCases := []struct {
		name   string
		setup  func(repo repository.UserPort) (*entity.User, error)
		errors bool
	}{
		{
			name: "happy path",
			setup: func(repo repository.UserPort) (*entity.User, error) {
				userInsert, err := repo.Insert(ctx, user)
				if err != nil {
					return nil, err
				}
				return userInsert, nil
			},
			errors: false,
		},
	}
}
