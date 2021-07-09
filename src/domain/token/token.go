package token

import (
	"fmt"
	"github.com/dbielecki97/bookstore-utils-go/crypto"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"strings"
	"time"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type Repository interface {
	GetById(string) (*Token, errs.RestErr)
	Create(*Token) errs.RestErr
	UpdateExpirationTime(Token) errs.RestErr
}

type Request struct {
	GrantType string `json:"grant_type,omitempty"`
	Scope     string `json:"scope"`
	// Used for password grant_type
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	// Used for client_credentials grand_type
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func (r Request) Validate() errs.RestErr {
	switch r.GrantType {
	case grantTypePassword:
		return nil
	case grantTypeClientCredentials:
		return nil
	default:
		return errs.NewBadRequestErr("invalid grant type")
	}
}

type Token struct {
	ID       string `json:"id,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
	ClientId int64  `json:"client_id,omitempty"`
	Expires  int64  `json:"expires,omitempty"`
}

func GetNewAccessToken(userId int64) Token {
	return Token{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
		UserId:  userId,
	}
}

func (t *Token) Generate() {
	t.ID = crypto.GetMd5(fmt.Sprintf("at-%d-%d-ran", t.UserId, t.Expires))
}

func (t *Token) IsExpired() bool {
	return time.Unix(t.Expires, 0).Before(time.Now().UTC())
}

func (t *Token) Validate() errs.RestErr {
	t.ID = strings.TrimSpace(t.ID)

	if t.ID == "" {
		return errs.NewBadRequestErr("invalid token id")
	}
	if t.UserId <= 0 {
		return errs.NewBadRequestErr("invalid user id")
	}
	if t.ClientId <= 0 {
		return errs.NewBadRequestErr("invalid client id")
	}
	if t.Expires <= 0 {
		return errs.NewBadRequestErr("invalid expiration time")
	}

	return nil
}
