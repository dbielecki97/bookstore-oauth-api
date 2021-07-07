package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime)
	assert.EqualValues(t, "password", grantTypePassword)
	assert.EqualValues(t, "client_credentials", grantTypeClientCredentials)
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken(0)

	assert.False(t, at.IsExpired(), "brand new access token should not be expired")
	assert.EqualValues(t, "", at.ID, "new access token should not have defined access token id")
	assert.True(t, at.UserId == 0, "new access token should not have associated user id")
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := Token{}

	assert.True(t, at.IsExpired(), "empty access token should be expired by default")

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "access token expiring in 3 hours should not be expired")
}
