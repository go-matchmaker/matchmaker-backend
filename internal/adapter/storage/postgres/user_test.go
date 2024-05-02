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

func getUserRepo(newDB db.EngineMaker) repository.UserMaker {
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
		setup  func(repo repository.UserMaker) error
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
	time.Sleep(5 * time.Second)
	require.NotNil(t, engine)
	userRepo := getUserRepo(engine)
	require.NotNil(t, userRepo)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
		setup  func(repo repository.UserMaker) (*uuid.UUID, error)
		input  *entity.User
		errors bool
	}{
		{
			name: "happy path",
			setup: func(repo repository.UserMaker) (*uuid.UUID, error) {
				id, err := repo.Insert(ctx, user)
				return id, err
			},
			input:  user,
			errors: false,
		},
		{
			name: "not found",
			setup: func(repo repository.UserMaker) (*uuid.UUID, error) {
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
			id, err := tc.setup(userRepo)
			require.NotNil(t, id)
			require.NoError(t, err)
			fmt.Println("ID:", *id)
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
		setup  func(repo repository.UserMaker) (*uuid.UUID, error)
		input  *entity.User
		errors bool
	}{
		{
			name: "happy path",
			setup: func(repo repository.UserMaker) (*uuid.UUID, error) {
				id, err := repo.Insert(ctx, user)
				return id, err
			},
			input:  user,
			errors: false,
		},
		{
			name: "not found",
			setup: func(repo repository.UserMaker) (*uuid.UUID, error) {
				randomUUID := uuid.New()
				return &randomUUID, nil
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
			id, err := tc.setup(userRepo)
			require.NotNil(t, id)
			require.NoError(t, err)
			fmt.Println("ID:", *id)
			err = userRepo.DeleteUser(ctx, *id)
			if tc.errors {
				fmt.Println("Error:", err)
				require.NotNil(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
	fmt.Println("Test done")
}
