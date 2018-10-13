package routing

import (
	"log"
	"net/http"

	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/radius"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/twilio"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
	"github.com/labstack/echo"
)

type CryptoRouter struct {
	worldbitService    worldbit.Worldbit
	radiusService      radius.Radius
	twilioService      twilio.Twilio
	pricingPlanService pricing_plans.PricingPlans
}

func NewCryptoRouter(
	worldbitService worldbit.Worldbit,
	radiusService radius.Radius,
	twilioService twilio.Twilio,
	pricingPlanService pricing_plans.PricingPlans,
) CryptoRouter {
	return CryptoRouter{
		worldbitService:    worldbitService,
		radiusService:      radiusService,
		twilioService:      twilioService,
		pricingPlanService: pricingPlanService,
	}
}

func (router CryptoRouter) Register(group *echo.Group) {
	group.GET("/plans", router.requestPlans)
	group.POST("/payment", router.requestPayment)
}

func (router CryptoRouter) requestPayment(context echo.Context) error {
	paymentRequest := new(PaymentRequest)

	if err := context.Bind(paymentRequest); err != nil {
		return err
	}

	plan, err := router.pricingPlanService.ByID(paymentRequest.PricingPlanID)
	if err != nil {
		return err
	}

	account, err := router.worldbitService.CreateAccount()
	if err != nil {
		return err
	}

	rate, err := router.worldbitService.GetExchangeRate()
	if err != nil {
		return err
	}

	exchange, err := router.worldbitService.CreateExchange(worldbit.CreateExchangeRequest{
		Amount:         plan.AmountUSD * rate,
		SenderCurrency: paymentRequest.Currency,
		Address:        account.Address,
	})
	if err != nil {
		return err
	}

	go func() {
		if err := router.worldbitService.MonitorExchangeStatus(exchange.StatusURL); err != nil {
			log.Printf("Monitoring exchange status exited with error: %s\n", err)
			return
		}

		usernamePassword, err := router.radiusService.CreateCredentials(
			radius.PricingPlan{
				Duration:  plan.Duration,
				MaxUsers:  plan.MaxUsers,
				UpLimit:   plan.UpLimit,
				DownLimit: plan.DownLimit,
				PurgeDays: plan.PurgeDays,
			},
		)
		if err != nil {
			log.Printf("Unable to generate credentials for radius. Error: %s\n", err)
			return
		}

		if err := router.twilioService.SendVoucher(paymentRequest.PhoneNumber, usernamePassword); err != nil {
			log.Printf("Unable to send voucher to phone number %s. Error: %s\n", paymentRequest.PhoneNumber, err)
			return
		}

		log.Printf("Successfully generated voucher %s for user with phone number %s and network id %s\n",
			usernamePassword,
			paymentRequest.PhoneNumber,
			paymentRequest.NetworkID,
		)
	}()

	return context.JSON(http.StatusOK, PaymentResponse{
		Address: exchange.Address,
		Amount:  exchange.Amount,
	})
}

func (router CryptoRouter) requestPlans(context echo.Context) error {
	plans, err := router.pricingPlanService.All()
	if err != nil {
		return err
	}
	return context.JSON(http.StatusOK, plans)
}
