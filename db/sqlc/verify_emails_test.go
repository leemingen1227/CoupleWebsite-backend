package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/leemingen1227/couple-server/util"
)

func TestCreateVerifyEmail (t *testing.T){
	user := createRandomUser(t)

	arg := CreateVerifyEmailParams {
		UserID: user.ID,
		Email: user.Email,
		SecretCode: util.RandomString(32),
	}

	verifyEmail, err := testQueries.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, verifyEmail)

	require.Equal(t, arg.UserID, verifyEmail.UserID)
	require.Equal(t, arg.Email, verifyEmail.Email)
	require.Equal(t, arg.SecretCode, verifyEmail.SecretCode)
}