package token

import "github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"

type Service interface {
	GetTokenById(string) (*AccessToken, *errors.RestErr)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTokenById(id string) (*AccessToken, *errors.RestErr) {
	return s.repo.GetById(id)
}
