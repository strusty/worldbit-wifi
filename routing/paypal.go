package routing

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/strusty/worldbit-wifi/services/paypal"
	"github.com/strusty/worldbit-wifi/services/pricing_plans"
	"github.com/strusty/worldbit-wifi/services/radius"
	"github.com/strusty/worldbit-wifi/services/twilio"
)

type PayPalRouter struct {
	payPalService       paypal.PayPal
	radiusService       radius.Radius
	twilioService       twilio.Twilio
	pricingPlansService pricing_plans.PricingPlans
}

func NewPayPalRouter(
	payPalService paypal.PayPal,
	radiusService radius.Radius,
	twilioService twilio.Twilio,
	pricingPlansService pricing_plans.PricingPlans,
) PayPalRouter {
	return PayPalRouter{
		payPalService:       payPalService,
		radiusService:       radiusService,
		twilioService:       twilioService,
		pricingPlansService: pricingPlansService,
	}
}

func (router PayPalRouter) Register(group *echo.Group) {
	group.POST("/payment", router.requestVoucher)
}

func (router PayPalRouter) requestVoucher(context echo.Context) error {
	request := new(PayPalVoucherRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	plan, err := router.pricingPlansService.ByID(request.PricingPlanID)
	if err != nil {
		return err
	}

	if err := router.payPalService.CheckSale(request.SaleID); err != nil {
		return err
	}

	voucher, err := router.radiusService.CreateCredentials(radius.PricingPlan{
		Duration:  plan.Duration,
		MaxUsers:  plan.MaxUsers,
		DownLimit: plan.DownLimit,
		PurgeDays: plan.PurgeDays,
	})
	if err != nil {
		return err
	}

	if err := router.twilioService.SendVoucher(request.PhoneNumber, voucher); err != nil {
		return err
	}

	if err := router.payPalService.PersistSale(request.SaleID, voucher); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, map[string]string{
		"voucher": voucher,
	})
}
