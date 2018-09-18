package routing

import (
	"testing"

	"git.sfxdx.ru/crystalline/wi-fi-backend/services/auth"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
)

func generateContext() echo.Context {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/test/10",
		strings.NewReader(`{}`),
	)

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return e.NewContext(request, recorder)
}

func generateContextWithInvalidBody() echo.Context {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/test/10",
		strings.NewReader(`{`),
	)

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return e.NewContext(request, recorder)
}

type AuthServiceMock struct {
	CreateCodeFn func(request auth.SendCodeRequest) (string, error)
}

func (mock AuthServiceMock) CreateCode(request auth.SendCodeRequest) (string, error) {
	return mock.CreateCodeFn(request)
}

type TwilioServiceMock struct {
	SendConfirmationCodeFn func(phoneNumber string, confirmationCode string) error
}

func (mock TwilioServiceMock) SendConfirmationCode(phoneNumber string, confirmationCode string) error {
	return mock.SendConfirmationCodeFn(phoneNumber, confirmationCode)
}

func TestNewAuthRouter(t *testing.T) {
	authService := &AuthServiceMock{}
	twilioService := &TwilioServiceMock{}

	testRouter := AuthRouter{
		authService:   authService,
		twilioService: twilioService,
	}

	router := NewAuthRouter(
		authService,
		twilioService,
	)

	assert.Equal(t, testRouter, router)
}

func TestAuthRouter_Register(t *testing.T) {
	AuthRouter{}.Register(echo.New().Group("/test"))
}

func TestAuthRouter_sendCode(t *testing.T) {
	testRouter := AuthRouter{
		authService: AuthServiceMock{
			CreateCodeFn: func(request auth.SendCodeRequest) (string, error) {
				return "", nil
			},
		},
		twilioService: TwilioServiceMock{
			SendConfirmationCodeFn: func(phoneNumber string, confirmationCode string) error {
				return nil
			},
		},
	}

	assert.NoError(t, testRouter.sendCode(generateContext()))

	testRouter.twilioService = TwilioServiceMock{
		SendConfirmationCodeFn: func(phoneNumber string, confirmationCode string) error {
			return errors.New("test_error")
		},
	}

	assert.Error(t, testRouter.sendCode(generateContext()))

	testRouter.authService = AuthServiceMock{
		CreateCodeFn: func(request auth.SendCodeRequest) (string, error) {
			return "", errors.New("test_error")
		},
	}

	assert.Error(t, testRouter.sendCode(generateContext()))
	assert.Error(t, testRouter.sendCode(generateContextWithInvalidBody()))

}

func TestAuthRouter_authenticate(t *testing.T) {

}
