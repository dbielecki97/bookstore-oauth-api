package rest

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/user"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
	"testing"
)

func setup() {
	httpmock.ActivateNonDefault(restClient.GetClient())
}

func shutdown() {
	httpmock.DeactivateAndReset()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

type timeoutError struct{}

func (e *timeoutError) Error() string { return "timeout error" }
func (e *timeoutError) Timeout() bool { return true }

func TestLoginUserTimeoutFromApi(t *testing.T) {
	httpmock.Reset()
	responder := httpmock.NewErrorResponder(&timeoutError{})

	httpmock.RegisterResponder(http.MethodPost, "/users/login", responder)

	repository := usersRepository{}

	_, apiErr := repository.LoginUser("email@gmail.com", "password")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.StatusCode())
	assert.EqualValues(t, "could not login user", apiErr.Message())
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	httpmock.Reset()
	responder, err := httpmock.NewJsonResponder(http.StatusInternalServerError, map[string]interface{}{"error": 123})
	if err != nil {
		t.Error("could not create responder")
	}
	httpmock.RegisterResponder(http.MethodPost, "/users/login", responder)

	repository := usersRepository{}

	_, restErr := repository.LoginUser("email@gmail.com", "password")
	if restErr.StatusCode() != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", restErr.StatusCode())
	}
	if !strings.Contains(restErr.Message(), "invalid error interface when trying to login user") {
		t.Errorf("invalid error has been thrown")
	}
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	httpmock.Reset()
	apiErr := errs.NewRestErr("invalid credentials", http.StatusUnauthorized, "unauthorized", nil)

	responder, err := httpmock.NewJsonResponder(http.StatusUnauthorized, apiErr)
	if err != nil {
		t.Error("could not create responder")
	}
	httpmock.RegisterResponder(http.MethodPost, "/users/login", responder)

	repository := usersRepository{}

	_, restErr := repository.LoginUser("email@gmail.com", "password")
	if restErr.StatusCode() != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", restErr.StatusCode())
	}
	if !strings.Contains(restErr.Message(), "invalid credentials") {
		t.Errorf("invalid error message")
	}
	if !strings.Contains(restErr.Err(), "unauthorized") {
		t.Errorf("invalid error")
	}
}

func TestLoginUserJsonResponseError(t *testing.T) {
	httpmock.Reset()

	responder, err := httpmock.NewJsonResponder(http.StatusOK, map[string]interface{}{"first_name": 123})
	if err != nil {
		t.Error("could not create responder")
	}
	httpmock.RegisterResponder(http.MethodPost, "/users/login", responder)

	repository := usersRepository{}

	_, restErr := repository.LoginUser("email@gmail.com", "password")
	if !strings.Contains(restErr.Message(), "error when trying to unmarshal users response") {
		t.Errorf("expected an marshalling error, got %v", restErr)
	}
}

func TestLoginUserNoError(t *testing.T) {
	httpmock.Reset()
	mockUser := user.User{
		ID:        12,
		FirstName: "Dawid",
		LastName:  "Bielecki",
		Email:     "email@gmail.com",
	}
	responder, err := httpmock.NewJsonResponder(http.StatusOK, mockUser)
	if err != nil {
		t.Error("could not create responder")
	}
	httpmock.RegisterResponder(http.MethodPost, "/users/login", responder)

	repository := usersRepository{}

	u, restErr := repository.LoginUser("email@gmail.com", "password")
	if restErr != nil {
		t.Errorf("expected no error, got %v", restErr)
	}

	assert.EqualValues(t, *u, mockUser)
}
