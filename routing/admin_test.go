package routing

import (
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"

	"git.sfxdx.ru/crystalline/wi-fi-backend/jwt"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/admins"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	userJSONCorrect = `{
		"Login": "admin",
		"Password": "admin"
	}`

	userJSONCorrectMixedCaseKeys = `{
		"lOgIn": "admin",
		"PassWord": "admin"
	}`

	userJSONWrongLogin = `{
		"Login": "wrong_login",
		"Password": "admin"
	}`

	userJSONWrongPassword = `{
		"Login": "admin",
		"Password": "wrong_password"
	}`

	jwtResponseCorrect = `{"token":"token"}`

	pricingPlanCreateCorrect = `{
		"amountUSD": 1,
        "duration": 1,
        "maxUsers": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	pricingPlanCreateCorrectResponse = `{
		"id": "id",
		"amountUSD": 1,
        "duration": 1,
        "maxUsers": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	pricingPlanCreateMissingField = `{
		"amountUSD": 1,
        "duration": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	pricingPlanUpdateCorrect = `{
		"id": "id",
		"amountUSD": 1,
        "duration": 1,
        "maxUsers": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	pricingPlanUpdateMissingID = `{
		"amountUSD": 1,
        "duration": 1,
        "maxUsers": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	pricingPlanUpdateMissingField = `{
		"id": "id",
        "duration": 1,
        "maxUsers": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	pricingPlanUpdateCorrectResponse = `{
		"id": "id",
		"amountUSD": 1,
        "duration": 1,
        "maxUsers": 1,
        "upLimit": 1,
        "downLimit": 1,
        "purgeDays": 0
	}`

	changePasswordRequestCorrect = `{
		"oldPassword": "admin",
		"newPassword": "newpass"
	}`

	changePasswordRequestWrongOldPassword = `{
		"oldPassword": "wrong_pass",
		"newPassword": "newpass"
	}`

	changePasswordRequestEmptyNewPassword = `{
		"oldPassword": "admin",
		"newPassword": "  "
	}`
)

func StripWhitespaces(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input)
}

func TestNewAdminRouter(t *testing.T) {
	adminService := &AdminServiceMock{}
	pricingPlanService := &PricingPlanServiceMock{}

	testRouter := AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlanService,
	}

	router := NewAdminRouter(
		adminService,
		pricingPlanService,
	)

	assert.Equal(t, testRouter, router)
}

func TestAdminRouter_Register(t *testing.T) {
	AdminRouter{}.Register(echo.New().Group("/test"))
}

func TestAdminRouter_login(t *testing.T) {
	adminService := &AdminServiceMock{
		LoginFn: func(request admins.LoginRequest) (*admins.JWTResponse, error) {
			if request.Login == "admin" && request.Password == "admin" {
				return &admins.JWTResponse{
					Token: "token",
				}, nil
			}
			return nil, errors.New("Login failed")
		},
	}
	pricingPlanService := &PricingPlanServiceMock{}
	e := echo.New()

	router := AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlanService,
	}

	t.Run("Correct", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(userJSONCorrect))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		if assert.NoError(t, router.login(context)) {
			assert.Equal(t, jwtResponseCorrect, rec.Body.String())
		}
	})

	t.Run("Mixed case keys", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(userJSONCorrectMixedCaseKeys))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		if assert.NoError(t, router.login(context)) {
			assert.Equal(t, jwtResponseCorrect, rec.Body.String())
		}
	})

	t.Run("Wrong login", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(userJSONWrongLogin))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		assert.Error(t, router.login(context))
	})

	t.Run("Wrong password", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(userJSONWrongLogin))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		assert.Error(t, router.login(context))
	})
}

func TestAdminRouter_changePassword(t *testing.T) {
	jwt.SetRandomSecret()
	adminService := &AdminServiceMock{
		ChangePasswordFn: func(id string, request admins.ChangePasswordRequest) error {
			if id == "id" && request.OldPassword == "admin" {
				return nil
			}
			return errors.New("Change password failed")
		},
	}
	pricingPlanService := &PricingPlanServiceMock{}

	e := echo.New()

	router := AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlanService,
	}

	t.Run("Correct", func(t *testing.T) {
		token, err := jwt.GenerateJWT("id")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		req := httptest.NewRequest(echo.PUT, "/password", strings.NewReader(changePasswordRequestCorrect))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(
			echo.HeaderAuthorization,
			strings.Join([]string{middleware.DefaultJWTConfig.AuthScheme, token}, " "),
		)
		jwtMiddleware, err := jwt.Middleware(middleware.DefaultSkipper)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		middlewaredChangePasswordRouter := jwtMiddleware(router.changePassword)

		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)

		assert.NoError(t, middlewaredChangePasswordRouter(context))
	})

	t.Run("Wrong id in jwt header", func(t *testing.T) {
		token, err := jwt.GenerateJWT("wrong_id")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		req := httptest.NewRequest(echo.PUT, "/password", strings.NewReader(changePasswordRequestCorrect))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(
			echo.HeaderAuthorization,
			strings.Join([]string{middleware.DefaultJWTConfig.AuthScheme, token}, " "),
		)
		jwtMiddleware, err := jwt.Middleware(middleware.DefaultSkipper)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		middlewaredChangePasswordRouter := jwtMiddleware(router.changePassword)

		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)

		assert.Error(t, middlewaredChangePasswordRouter(context))
	})

	t.Run("Malformed jwt header", func(t *testing.T) {
		token := "abcd"

		req := httptest.NewRequest(echo.PUT, "/password", strings.NewReader(changePasswordRequestCorrect))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(
			echo.HeaderAuthorization,
			strings.Join([]string{middleware.DefaultJWTConfig.AuthScheme, token}, " "),
		)
		jwtMiddleware, err := jwt.Middleware(middleware.DefaultSkipper)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		middlewaredChangePasswordRouter := jwtMiddleware(router.changePassword)

		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)

		assert.Error(t, middlewaredChangePasswordRouter(context))
	})

	t.Run("Wrong old password", func(t *testing.T) {
		token, err := jwt.GenerateJWT("id")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		req := httptest.NewRequest(echo.PUT, "/password", strings.NewReader(changePasswordRequestWrongOldPassword))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(
			echo.HeaderAuthorization,
			strings.Join([]string{middleware.DefaultJWTConfig.AuthScheme, token}, " "),
		)
		jwtMiddleware, err := jwt.Middleware(middleware.DefaultSkipper)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		middlewaredChangePasswordRouter := jwtMiddleware(router.changePassword)

		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)

		assert.Error(t, middlewaredChangePasswordRouter(context))
	})

	t.Run("Empty new password", func(t *testing.T) {
		token, err := jwt.GenerateJWT("id")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		req := httptest.NewRequest(echo.PUT, "/password", strings.NewReader(changePasswordRequestEmptyNewPassword))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(
			echo.HeaderAuthorization,
			strings.Join([]string{middleware.DefaultJWTConfig.AuthScheme, token}, " "),
		)
		jwtMiddleware, err := jwt.Middleware(middleware.DefaultSkipper)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		middlewaredChangePasswordRouter := jwtMiddleware(router.changePassword)

		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)

		assert.Error(t, middlewaredChangePasswordRouter(context))
	})
}

func TestAdminRouter_createPlan(t *testing.T) {
	adminService := &AdminServiceMock{}
	pricingPlanService := &PricingPlanServiceMock{
		CreateFn: func(plan *pricing_plans.PricingPlan) error {
			plan.ID = "id"
			return nil
		},
	}
	e := echo.New()

	router := AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlanService,
	}

	t.Run("Correct", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/plans", strings.NewReader(pricingPlanCreateCorrect))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		if assert.NoError(t, router.createPlan(context)) {
			assert.Equal(t, StripWhitespaces(pricingPlanCreateCorrectResponse), rec.Body.String())
		}
	})

	t.Run("Missing field", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/plans", strings.NewReader(pricingPlanCreateMissingField))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		assert.Error(t, router.createPlan(context))
	})
}

func TestAdminRouter_updatePlan(t *testing.T) {
	adminService := &AdminServiceMock{}
	pricingPlanService := &PricingPlanServiceMock{
		UpdateFn: func(plan *pricing_plans.PricingPlan) error {
			return nil
		},
	}
	e := echo.New()

	router := AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlanService,
	}

	t.Run("Correct", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/plans/id", strings.NewReader(pricingPlanUpdateCorrect))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		context.SetParamNames("id")
		context.SetParamValues("id")
		if assert.NoError(t, router.updatePlan(context)) {
			assert.Equal(t, StripWhitespaces(pricingPlanUpdateCorrectResponse), rec.Body.String())
		}
	})

	t.Run("Missing id", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/plans/id", strings.NewReader(pricingPlanUpdateMissingID))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		context.SetParamNames("id")
		context.SetParamValues("id")
		assert.Error(t, router.updatePlan(context))
	})

	t.Run("Missing field", func(t *testing.T) {
		req := httptest.NewRequest(echo.PUT, "/plans/id", strings.NewReader(pricingPlanUpdateMissingField))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		context.SetParamNames("id")
		context.SetParamValues("id")
		assert.Error(t, router.updatePlan(context))
	})
}

func TestAdminRouter_deletePlan(t *testing.T) {
	adminService := &AdminServiceMock{}
	pricingPlanService := &PricingPlanServiceMock{
		DeleteFn: func(id string) error {
			return nil
		},
	}
	e := echo.New()

	router := AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlanService,
	}

	t.Run("Correct", func(t *testing.T) {
		req := httptest.NewRequest(echo.DELETE, "/plans/id", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)
		context.SetParamNames("id")
		context.SetParamValues("id")
		assert.NoError(t, router.deletePlan(context))
	})
}
