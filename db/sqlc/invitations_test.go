package db

import (
	"context"
	"database/sql"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/leemingen1227/couple-server/util"
)

func TestCreateInvitation(t *testing.T){
	user := createRandomUser(t)
	var err error
	//Update the user's email to verified
	user, err = testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: user.ID,
		IsEmailVerified: sql.NullBool{
			Bool: true,
			Valid: true,
		},
	})

	arg := CreateInvitationParams{
		InviterID: user.ID,
		InviteeEmail: user.Email,
		InvitationToken: util.RandomString(32),
	}

	invitation, err := testQueries.CreateInvitation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invitation)

	require.Equal(t, arg.InviterID, invitation.InviterID)
	require.Equal(t, arg.InviteeEmail, invitation.InviteeEmail)
	require.Equal(t, false, invitation.IsAccepted)
}