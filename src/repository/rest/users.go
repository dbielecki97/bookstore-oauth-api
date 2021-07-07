package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/user"
	"github.com/dbielecki97/bookstore-oauth-api/src/utils/errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type UsersRepository interface {
	LoginUser(string, string) (*user.User, *errors.RestErr)
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

func (u usersRepository) LoginUser(email string, password string) (*user.User, *errors.RestErr) {
	lr := user.LoginRequest{
		Email:    email,
		Password: password,
	}

	res, err := restClient.R().SetBody(lr).Post("/users/login")
	if err != nil {
		fmt.Println(err)
		return nil, errors.NewInternalServerError("restclient error")
	}

	if res.StatusCode() > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(res.Body(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var usr user.User
	err = json.Unmarshal(res.Body(), &usr)
	if err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users response")
	}

	return &usr, nil
}
