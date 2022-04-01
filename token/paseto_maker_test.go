package token

import (
	"testing"
	"time"

	"github.com/hhong0326/goPostgresqlDocker.git/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

// Same interface with JWTMaker
func TestPasetoMaker(t *testing.T) {

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute // 1 minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())

	require.Nil(t, payload)
}

// require correction
func TestInvalidPasetoTokenAlg(t *testing.T) {

	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	token, err := paseto.NewV2().Encrypt([]byte("12345678901234567890123456789032"), payload, nil)

	t.Log(token)
	maker, err := NewPasetoMaker(token)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
