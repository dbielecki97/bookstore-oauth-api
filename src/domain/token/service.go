package token

import (
	"errors"
	"github.com/dbielecki97/bookstore-oauth-api/src/repository/rest"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"strings"
)

type Service interface {
	GetTokenById(string) (*Token, *errs.RestErr)
	CreateToken(Request) (*Token, *errs.RestErr)
	UpdateTokenExpiration(Token) *errs.RestErr
}

type service struct {
	tokenRepo Repository
	restRepo  rest.UsersRepository
}

func NewService(tokenRepo Repository, restRepo rest.UsersRepository) Service {
	return &service{tokenRepo: tokenRepo, restRepo: restRepo}
}

func (s *service) GetTokenById(id string) (*Token, *errs.RestErr) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errs.NewBadRequestErr("invalid token id")
	}

	return s.tokenRepo.GetById(id)
}

func (s *service) CreateToken(r Request) (*Token, *errs.RestErr) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	var token *Token
	var err *errs.RestErr

	switch r.GrantType {
	case grantTypePassword:
		token, err = s.generateTokenFromPassword(r)
	case grantTypeClientCredentials:
		return nil, errs.NewInternalServerErr("error processing request", errors.New("not implemented"))
	}
	if err != nil {
		return nil, err
	}

	err = s.tokenRepo.Create(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *service) generateTokenFromPassword(r Request) (*Token, *errs.RestErr) {
	user, err := s.restRepo.LoginUser(r.Username, r.Password)
	if err != nil {
		return nil, err
	}

	token := GetNewAccessToken(user.ID)
	token.Generate()
	return &token, nil
}

func (s *service) UpdateTokenExpiration(t Token) *errs.RestErr {
	if err := t.Validate(); err != nil {
		return err
	}

	return s.tokenRepo.UpdateExpirationTime(t)
}
