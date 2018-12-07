package routing

import (
	"testing"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/services/pricing_plans"
)

func TestNewPricingPlansRouter(t *testing.T) {
	pricingPlansServiceMock := &PricingPlanServiceMock{
		AllFn: func() ([]pricing_plans.PricingPlan, error) {
			return []pricing_plans.PricingPlan{}, nil
		},
	}

	testRouter := PricingPlansRouter{
		pricingPlansService: pricingPlansServiceMock,
	}
	router := NewPricingPlansRouter(pricingPlansServiceMock)

	t.Run("Initialization", func(t *testing.T) {
		if assert.Equal(t, testRouter, router) {
			router.Register(echo.New().Group("/test"))
			t.Run("Success", func(t *testing.T) {
				assert.NoError(t, router.plans(generateContext()))
			})

			pricingPlansServiceMock.AllFn = func() ([]pricing_plans.PricingPlan, error) {
				return nil, errors.New("test_error")
			}

			t.Run("Error", func(t *testing.T) {
				assert.Error(t, router.plans(generateContext()))
			})
		}
	})
}
