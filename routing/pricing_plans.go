package routing

import (
	"net/http"

	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"github.com/labstack/echo"
)

type PricingPlansRouter struct {
	pricingPlansService pricing_plans.PricingPlans
}

func NewPricingPlansRouter(pricingPlansService pricing_plans.PricingPlans) PricingPlansRouter {
	return PricingPlansRouter{
		pricingPlansService: pricingPlansService,
	}
}

func (router PricingPlansRouter) Register(group *echo.Group) {
	group.GET("", router.plans)
}

func (router PricingPlansRouter) plans(context echo.Context) error {
	plans, err := router.pricingPlansService.All()
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, plans)
}
