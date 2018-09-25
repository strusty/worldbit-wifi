package routing

import (
	"testing"

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
	CreateCodeFn func(phoneNumber string) (string, error)
	VerifyCodeFn func(confirmationCode string) error
}

func (mock AuthServiceMock) CreateCode(phoneNumber string) (string, error) {
	return mock.CreateCodeFn(phoneNumber)
}

func (mock AuthServiceMock) VerifyCode(confirmationCode string) error {
	return mock.VerifyCodeFn(confirmationCode)
}

type TwilioServiceMock struct {
	SendConfirmationCodeFn func(phoneNumber string, confirmationCode string) error
}

func (mock TwilioServiceMock) SendConfirmationCode(phoneNumber string, confirmationCode string) error {
	return mock.SendConfirmationCodeFn(phoneNumber, confirmationCode)
}

type CaptchaServiceMock struct {
	CheckCaptchaFn func(responseToken string) (bool, error)
}

func (mock CaptchaServiceMock) CheckCaptcha(responseToken string) (bool, error) {
	return mock.CheckCaptchaFn(responseToken)
}

func TestNewAuthRouter(t *testing.T) {
	authService := &AuthServiceMock{}
	twilioService := &TwilioServiceMock{}
	captchaService := &CaptchaServiceMock{}

	testRouter := AuthRouter{
		authService:    authService,
		captchaService: captchaService,
		twilioService:  twilioService,
	}

	router := NewAuthRouter(
		authService,
		captchaService,
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
			CreateCodeFn: func(phoneNumber string) (string, error) {
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
		CreateCodeFn: func(phoneNumber string) (string, error) {
			return "", errors.New("test_error")
		},
	}

	assert.Error(t, testRouter.sendCode(generateContext()))
	assert.Error(t, testRouter.sendCode(generateContextWithInvalidBody()))

}

func TestAuthRouter_authenticate(t *testing.T) {
	testRouter := AuthRouter{
		authService: AuthServiceMock{
			VerifyCodeFn: func(confirmationCode string) error {
				return nil
			},
		},
		captchaService: CaptchaServiceMock{
			CheckCaptchaFn: func(responseToken string) (bool, error) {
				return true, nil
			},
		},
	}

	assert.NoError(t, testRouter.authenticate(generateContext()))

	testRouter.authService = AuthServiceMock{
		VerifyCodeFn: func(confirmationCode string) error {
			return errors.New("test_error")
		},
	}

	assert.Error(t, testRouter.authenticate(generateContext()))

	testRouter.captchaService = CaptchaServiceMock{
		CheckCaptchaFn: func(responseToken string) (bool, error) {
			return false, nil
		},
	}

	assert.Error(t, testRouter.authenticate(generateContext()))

	testRouter.captchaService = CaptchaServiceMock{
		CheckCaptchaFn: func(responseToken string) (bool, error) {
			return false, errors.New("test_error")
		},
	}

	assert.Error(t, testRouter.authenticate(generateContext()))
	assert.Error(t, testRouter.authenticate(generateContextWithInvalidBody()))
}
