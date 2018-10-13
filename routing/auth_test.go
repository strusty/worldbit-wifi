package routing

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

func generateContextWithPricingPlanID(id string) echo.Context {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/test/10",
		strings.NewReader(`{ "PricingPlanID": "`+id+`" }`),
	)

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return e.NewContext(request, recorder)
}

func generateContextForPricingPlan(operation string) echo.Context {
	e := echo.New()

	recorder := httptest.NewRecorder()
	requestMethod := http.MethodGet
	requestPath := ""
	requestBody := ""

	switch operation {
	case "create":
		requestMethod = http.MethodPost
		requestPath = "/test"
		requestBody = `{ "amountUSD": 1 }`
	case "update":
		requestMethod = http.MethodPut
		requestPath = "/test"
		requestBody = `{ "ID": "id", "amountUSD": 2 }`
	case "delete":
		requestMethod = http.MethodDelete
		requestPath = "/test"
	}

	request := httptest.NewRequest(
		requestMethod,
		requestPath,
		strings.NewReader(requestBody),
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

func generateContextForLogin() echo.Context {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/test/10",
		strings.NewReader(`{ "Login": "admin", "Password": "password" }`),
	)

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return e.NewContext(request, recorder)
}

func generateContextForChangePassword(password string) echo.Context {
	e := echo.New()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPut,
		"/test/10",
		strings.NewReader(`{ "OldPassword": "password", "NewPassword": "`+password+`"  }`),
	)

	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return e.NewContext(request, recorder)
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
		twilioService: &TwilioServiceMock{
			SendConfirmationCodeFn: func(phoneNumber string, confirmationCode string) error {
				return nil
			},
		},
	}
	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, testRouter.sendCode(generateContext()))
	})

	testRouter.twilioService = &TwilioServiceMock{
		SendConfirmationCodeFn: func(phoneNumber string, confirmationCode string) error {
			return errors.New("test_error")
		},
	}

	t.Run("Unable to send confirmation code error", func(t *testing.T) {
		assert.Error(t, testRouter.sendCode(generateContext()))
	})

	testRouter.authService = AuthServiceMock{
		CreateCodeFn: func(phoneNumber string) (string, error) {
			return "", errors.New("test_error")
		},
	}

	t.Run("Unable to create confirmation code error", func(t *testing.T) {
		assert.Error(t, testRouter.sendCode(generateContext()))
	})

	t.Run("Invalid body error", func(t *testing.T) {
		assert.Error(t, testRouter.sendCode(generateContextWithInvalidBody()))
	})
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

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, testRouter.authenticate(generateContext()))
	})

	testRouter.authService = AuthServiceMock{
		VerifyCodeFn: func(confirmationCode string) error {
			return errors.New("test_error")
		},
	}

	t.Run("Code verification error", func(t *testing.T) {
		assert.Error(t, testRouter.authenticate(generateContext()))
	})

	testRouter.captchaService = CaptchaServiceMock{
		CheckCaptchaFn: func(responseToken string) (bool, error) {
			return false, nil
		},
	}

	t.Run("Check captcha failed", func(t *testing.T) {
		assert.Error(t, testRouter.authenticate(generateContext()))
	})

	testRouter.captchaService = CaptchaServiceMock{
		CheckCaptchaFn: func(responseToken string) (bool, error) {
			return false, errors.New("test_error")
		},
	}

	t.Run("Check captcha error", func(t *testing.T) {
		assert.Error(t, testRouter.authenticate(generateContext()))
	})

	t.Run("Invalid body error", func(t *testing.T) {
		assert.Error(t, testRouter.authenticate(generateContextWithInvalidBody()))
	})
}
