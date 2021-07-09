package cassandra

import (
	"fmt"
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"github.com/dbielecki97/bookstore-utils-go/logger"
	"github.com/gocql/gocql"
)

const (
	getToken             = "SELECT id, user_id, client_id, expires from tokens where id=?;"
	createToken          = "INSERT INTO tokens(id, user_id, client_id, expires) values(?, ?, ?, ?);"
	updateExpirationTime = "UPDATE tokens set expires = ? where id = ?;"
)

type cassandraRepository struct {
	session *gocql.Session
}

func New(session *gocql.Session) token.Repository {
	return &cassandraRepository{session: session}
}

func (r cassandraRepository) GetById(id string) (*token.Token, errs.RestErr) {
	var t token.Token
	if err := r.session.Query(getToken, id).Scan(&t.ID, &t.UserId, &t.ClientId, &t.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errs.NewNotFoundErr(fmt.Sprintf("token with id %s not found", id))
		}
		logger.Error("could not get token by id", err)
		return nil, errs.NewInternalServerErr("could not get token", err)
	}

	return &t, nil
}

func (r cassandraRepository) Create(t *token.Token) errs.RestErr {
	err := r.session.Query(createToken, t.ID, t.UserId, t.ClientId, t.Expires).Exec()
	if err != nil {
		logger.Error("could not create token", err)
		return errs.NewInternalServerErr("could not create token", err)
	}

	return nil
}

func (r cassandraRepository) UpdateExpirationTime(t token.Token) errs.RestErr {
	err := r.session.Query(updateExpirationTime, t.Expires, t.ID).Exec()
	if err != nil {
		logger.Error("could not set new expiration time", err)
		return errs.NewInternalServerErr("could not refresh token", err)
	}

	return nil
}
