package db

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"
)

type cassandraRepository struct{}

func NewRepo() token.Repository {
	return &cassandraRepository{}
}

func (r cassandraRepository) GetById(id string) (*token.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("database connection not implemented yet")
}
