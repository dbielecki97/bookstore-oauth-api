package token

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"
	"strings"
	"time"
)

const (
	expirationTime = 24
)

type Token struct {
	ID       string `json:"id,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
	ClientId int64  `json:"client_id,omitempty"`
	Expires  int64  `json:"expires,omitempty"`
}

type Repository interface {
	GetById(string) (*Token, *errors.RestErr)
	Create(Token) *errors.RestErr
	UpdateExpirationTime(Token) *errors.RestErr
}

func GetNewAccessToken() Token {
	return Token{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (t *Token) IsExpired() bool {
	return time.Unix(t.Expires, 0).Before(time.Now().UTC())
}

func (t *Token) Validate() *errors.RestErr {
	t.ID = strings.TrimSpace(t.ID)

	if t.ID == "" {
		return errors.NewBadRequestError("invalid token id")
	}
	if t.UserId <= 0 {
		return errors.NewBadRequestError("invalid user id")
	}
	if t.ClientId <= 0 {
		return errors.NewBadRequestError("invalid client id")
	}
	if t.Expires <= 0 {
		return errors.NewBadRequestError("invalid expiration time")
	}

	return nil
}
