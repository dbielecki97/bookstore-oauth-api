package token

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"
	"strings"
)

type Service interface {
	GetTokenById(string) (*Token, *errors.RestErr)
	CreateToken(Token) *errors.RestErr
	UpdateTokenExpiration(Token) *errors.RestErr
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTokenById(id string) (*Token, *errors.RestErr) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return nil, errors.NewBadRequestError("invalid token id")
	}

	return s.repo.GetById(id)
}

func (s *service) CreateToken(t Token) *errors.RestErr {
	if err := t.Validate(); err != nil {
		return err
	}

	return s.repo.Create(t)
}

func (s *service) UpdateTokenExpiration(t Token) *errors.RestErr {
	if err := t.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateExpirationTime(t)
}
