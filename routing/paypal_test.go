package routing

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/services/pricing_plans"
	"github.com/strusty/worldbit-wifi/services/radius"
)

func TestPayPalRouter(t *testing.T) {

	payPalServiceMock := &PayPalServiceMock{
		CheckSaleFn: func(saleID string) error {
			return nil
		},
		PersistSaleFn: func(saleID string, voucher string) error {
			return nil
		},
	}

	radiusServiceMock := &RadiusServiceMock{
		CreateCredentialsFn: func(plan radius.PricingPlan) (string, error) {
			return "voucher", nil
		},
	}

	twilioServiceMock := &TwilioServiceMock{
		SendVoucherFn: func(phoneNumber string, voucher string) error {
			return nil
		},
	}

	pricingPlanServiceMock := &PricingPlanServiceMock{
		ByIDFn: func(id string) (*pricing_plans.PricingPlan, error) {
			return &pricing_plans.PricingPlan{
				ID:        "id",
				AmountUSD: 10,
				Duration:  20,
				MaxUsers:  30,
				UpLimit:   40,
				DownLimit: 50,
				PurgeDays: 60,
			}, nil
		},
	}

	testRouter := PayPalRouter{
		payPalService:       payPalServiceMock,
		radiusService:       radiusServiceMock,
		twilioService:       twilioServiceMock,
		pricingPlansService: pricingPlanServiceMock,
	}

	router := NewPayPalRouter(
		payPalServiceMock,
		radiusServiceMock,
		twilioServiceMock,
		pricingPlanServiceMock,
	)

	t.Run("Initialization", func(t *testing.T) {
		if assert.Equal(t, testRouter, router) {
			router.Register(echo.New().Group("/test"))

			t.Run("requestVoucher", func(t *testing.T) {
				t.Run("Success", func(t *testing.T) {
					assert.NoError(t, router.requestVoucher(generateContext()))
				})

				payPalServiceMock.PersistSaleFn = func(saleID string, voucher string) error {
					return errors.New("test_error")
				}
				t.Run("Persist sale error", func(t *testing.T) {
					assert.Error(t, router.requestVoucher(generateContext()))
				})

				twilioServiceMock.SendVoucherFn = func(phoneNumber string, voucher string) error {
					return errors.New("test_error")
				}
				t.Run("Send voucher error", func(t *testing.T) {
					assert.Error(t, router.requestVoucher(generateContext()))
				})

				radiusServiceMock.CreateCredentialsFn = func(plan radius.PricingPlan) (string, error) {
					return "", errors.New("test_error")
				}
				t.Run("Create credentials error", func(t *testing.T) {
					assert.Error(t, router.requestVoucher(generateContext()))
				})

				payPalServiceMock.CheckSaleFn = func(saleID string) error {
					return errors.New("test_error")
				}
				t.Run("Check sale error", func(t *testing.T) {
					assert.Error(t, router.requestVoucher(generateContext()))
				})

				pricingPlanServiceMock.ByIDFn = func(id string) (*pricing_plans.PricingPlan, error) {
					return nil, errors.New("test_error")
				}
				t.Run("Pricing plan by id error", func(t *testing.T) {
					assert.Error(t, router.requestVoucher(generateContext()))
				})

				t.Run("Invalid json error", func(t *testing.T) {
					assert.Error(t, router.requestVoucher(generateContextWithInvalidBody()))
				})
			})
		}
	})
}
