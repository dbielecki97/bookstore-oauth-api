package token

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"
	"time"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
	ClientId    int64  `json:"client_id,omitempty"`
	Expires     int64  `json:"expires,omitempty"`
}

type Repository interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

func GetNewAccessToken() AccessToken {
	return AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (a *AccessToken) IsExpired() bool {
	return time.Unix(a.Expires, 0).Before(time.Now().UTC())
}
