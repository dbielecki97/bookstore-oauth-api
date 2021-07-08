package rest

import (
	"encoding/json"
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/user"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"github.com/dbielecki97/bookstore-utils-go/logger"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type UsersRepository interface {
	LoginUser(string, string) (*user.User, *errs.RestErr)
}

type usersRepository struct {
}

func New() UsersRepository {
	return &usersRepository{}
}

var (
	restClient = resty.NewWithClient(&http.Client{
		Timeout: 300 * time.Millisecond,
	}).SetHostURL("http://localhost:8080")
)

func (u usersRepository) LoginUser(email string, password string) (*user.User, *errs.RestErr) {
	lr := user.LoginRequest{
		Email:    email,
		Password: password,
	}

	res, err := restClient.R().SetBody(lr).Post("/users/login")
	if err != nil {
		logger.Error("POST request failed", err)
		return nil, errs.NewInternalServerErr("could not login user", err)
	}

	if res.StatusCode() > 299 {
		var restErr errs.RestErr
		err := json.Unmarshal(res.Body(), &restErr)
		if err != nil {
			return nil, errs.NewInternalServerErr("invalid error interface when trying to login user", err)
		}
		return nil, &restErr
	}

	var usr user.User
	err = json.Unmarshal(res.Body(), &usr)
	if err != nil {
		return nil, errs.NewInternalServerErr("error when trying to unmarshal users response", err)
	}

	return &usr, nil
}
