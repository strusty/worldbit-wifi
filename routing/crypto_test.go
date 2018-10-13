package routing

import (
	"testing"
	"time"

	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/radius"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewCryptoRouter(t *testing.T) {
	worldbitService := &WorldbitServiceMock{}
	radiusService := &RadiusServiceMock{}
	twilioService := &TwilioServiceMock{}
	pricingPlanService := &PricingPlanServiceMock{}

	testRouter := CryptoRouter{
		worldbitService:    worldbitService,
		radiusService:      radiusService,
		twilioService:      twilioService,
		pricingPlanService: pricingPlanService,
	}

	router := NewCryptoRouter(
		worldbitService,
		radiusService,
		twilioService,
		pricingPlanService,
	)

	assert.Equal(t, testRouter, router)
}

func TestCryptoRouter_Register(t *testing.T) {
	CryptoRouter{}.Register(echo.New().Group("/test"))
}

func TestCryptoRouter_requestPayment(t *testing.T) {
	worldbitService := &WorldbitServiceMock{
		CreateAccountFn: func() (*worldbit.CreateAccountResponseData, error) {
			return &worldbit.CreateAccountResponseData{}, nil
		},
		GetExchangeRateFn: func() (float64, error) {
			return 0, nil
		},
		CreateExchangeFn: func(request worldbit.CreateExchangeRequest) (*worldbit.CreateExchangeResult, error) {
			return &worldbit.CreateExchangeResult{}, nil
		},
		MonitorExchangeStatusFn: func(statusURL string) error {
			return nil
		},
	}

	radiusService := &RadiusServiceMock{
		CreateCredentialsFn: func(plan radius.PricingPlan) (string, error) {
			return "credentials", nil
		},
	}

	twilioService := &TwilioServiceMock{
		SendVoucherFn: func(phoneNumber string, voucher string) error {
			return nil
		},
	}

	pricingPlanService := &PricingPlanServiceMock{
		ByIDFn: func(id string) (*pricing_plans.PricingPlan, error) {
			if id == "id" {
				return &pricing_plans.PricingPlan{
					ID:        id,
					AmountUSD: 1,
					Duration:  1,
				}, nil
			}
			return nil, errors.New("error_bad_id")
		},
	}

	router := CryptoRouter{
		worldbitService:    worldbitService,
		radiusService:      radiusService,
		twilioService:      twilioService,
		pricingPlanService: pricingPlanService,
	}

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.True(t, radiusService.CreateCredentialsFnInvoked)
		assert.True(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.MonitorExchangeStatusFnInvoked = false
	radiusService.CreateCredentialsFnInvoked = false
	twilioService.SendVoucherFnInvoked = false
	twilioService.SendVoucherFn = func(phoneNumber string, voucher string) error {
		return errors.New("test_error")
	}

	t.Run("Send voucher error", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.True(t, radiusService.CreateCredentialsFnInvoked)
		assert.True(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.MonitorExchangeStatusFnInvoked = false
	radiusService.CreateCredentialsFnInvoked = false
	radiusService.CreateCredentialsFn = func(plan radius.PricingPlan) (string, error) {
		return "", errors.New("test_error")
	}
	twilioService.SendVoucherFnInvoked = false

	t.Run("Create voucher error", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.True(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.MonitorExchangeStatusFnInvoked = false
	worldbitService.MonitorExchangeStatusFn = func(statusURL string) error {
		return errors.New("test_error")
	}
	radiusService.CreateCredentialsFnInvoked = false

	t.Run("Monitor exchange status error", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.CreateExchangeFn = func(request worldbit.CreateExchangeRequest) (*worldbit.CreateExchangeResult, error) {
		return nil, errors.New("test_error")
	}
	worldbitService.MonitorExchangeStatusFnInvoked = false

	t.Run("Create exchange error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.GetExchangeRateFn = func() (float64, error) {
		return 0, errors.New("test_error")
	}

	t.Run("Get exchange rate error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.CreateAccountFn = func() (*worldbit.CreateAccountResponseData, error) {
		return nil, errors.New("test_error")
	}

	t.Run("Create account error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContextWithPricingPlanID("id")))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	t.Run("Invalid body error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContextWithInvalidBody()))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	t.Run("Pricing plan ID error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContextWithPricingPlanID("bad_id")))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, radiusService.CreateCredentialsFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

}

func TestCryptoRouter_requestPlans(t *testing.T) {
	worldbitService := &WorldbitServiceMock{
		CreateAccountFn: func() (*worldbit.CreateAccountResponseData, error) {
			return &worldbit.CreateAccountResponseData{}, nil
		},
		GetExchangeRateFn: func() (float64, error) {
			return 0, nil
		},
		CreateExchangeFn: func(request worldbit.CreateExchangeRequest) (*worldbit.CreateExchangeResult, error) {
			return &worldbit.CreateExchangeResult{}, nil
		},
		MonitorExchangeStatusFn: func(statusURL string) error {
			return nil
		},
	}

	radiusService := &RadiusServiceMock{
		CreateCredentialsFn: func(plan radius.PricingPlan) (string, error) {
			return "", nil
		},
	}

	twilioService := &TwilioServiceMock{
		SendVoucherFn: func(phoneNumber string, voucher string) error {
			return nil
		},
	}

	pricingPlanService := &PricingPlanServiceMock{
		AllFn: func() ([]pricing_plans.PricingPlan, error) {
			return []pricing_plans.PricingPlan{
				{
					ID:        "id",
					AmountUSD: 1,
				},
			}, nil
		},
	}

	router := CryptoRouter{
		worldbitService:    worldbitService,
		radiusService:      radiusService,
		twilioService:      twilioService,
		pricingPlanService: pricingPlanService,
	}

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, router.requestPlans(generateContext()))
	})

	pricingPlanService.AllFn = func() ([]pricing_plans.PricingPlan, error) {
		return nil, errors.New("test_error")
	}

	t.Run("Error", func(t *testing.T) {
		assert.Error(t, router.requestPlans(generateContext()))
	})
}
