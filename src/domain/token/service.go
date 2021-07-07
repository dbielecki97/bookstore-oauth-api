package token

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/repository/rest"
	"github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"
	"strings"
)

type Service interface {
	GetTokenById(string) (*Token, *errors.RestErr)
	CreateToken(Request) (*Token, *errors.RestErr)
	UpdateTokenExpiration(Token) *errors.RestErr
}

type service struct {
	tokenRepo Repository
	restRepo  rest.UsersRepository
}

func NewService(tokenRepo Repository, restRepo rest.UsersRepository) Service {
	return &service{tokenRepo: tokenRepo, restRepo: restRepo}
}

func (s *service) GetTokenById(id string) (*Token, *errors.RestErr) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.NewBadRequestError("invalid token id")
	}

	return s.tokenRepo.GetById(id)
}

func (s *service) CreateToken(r Request) (*Token, *errors.RestErr) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	var token *Token
	var err *errors.RestErr

	switch r.GrantType {
	case grantTypePassword:
		token, err = s.generateTokenFromPassword(r)
	case grantTypeClientCredentials:
		return nil, errors.NewInternalServerError("not implemented")
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

func (s *service) generateTokenFromPassword(r Request) (*Token, *errors.RestErr) {
	user, err := s.restRepo.LoginUser(r.Username, r.Password)
	if err != nil {
		return nil, err
	}

	token := GetNewAccessToken(user.ID)
	token.Generate()
	return &token, nil
}

func (s *service) UpdateTokenExpiration(t Token) *errors.RestErr {
	if err := t.Validate(); err != nil {
		return err
	}

	return s.tokenRepo.UpdateExpirationTime(t)
}
