package db

import (
	"context"
	"github.com/leemingen1227/couple-server/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User{
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	uid, err := uuid.NewRandom()
	arg := CreateUserParams{
		ID:             uid,
		Email:          util.RandomEmail(),
		Name:           util.RandomOwner(),
		PasswordDigest: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.PasswordDigest, user.PasswordDigest)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreateTime)
	require.NotZero(t, user.UpdateTime)
	
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
