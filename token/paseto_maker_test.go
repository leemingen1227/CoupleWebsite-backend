package token

import (
	"testing"
	"time"
	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func TestPasetoMaker(t *testing.T) {
    maker, err := NewPasetoMaker(util.RandomString(32))
    require.NoError(t, err)

    userID := uuid.New()
    duration := time.Minute

    issuedAt := time.Now()
    expiredAt := issuedAt.Add(duration)

    token, payload, err := maker.CreateToken(userID, duration)
    require.NoError(t, err)
    require.NotEmpty(t, token)
    require.NotEmpty(t, payload)

    payload, err = maker.VerifyToken(token)
    require.NoError(t, err)
    require.NotEmpty(t, token)

    require.NotZero(t, payload.ID)
    require.Equal(t, userID, payload.UserID)
    require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
    require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
    maker, err := NewPasetoMaker(util.RandomString(32))
    require.NoError(t, err)

	userID := uuid.New()

    token, payload, err := maker.CreateToken(userID, -time.Minute)
    require.NoError(t, err)
    require.NotEmpty(t, token)
    require.NotEmpty(t, payload)

    payload, err = maker.VerifyToken(token)
    require.Error(t, err)
    require.EqualError(t, err, ErrExpiredToken.Error())
    require.Nil(t, payload)
}