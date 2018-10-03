package routing

import (
	"log"
	"net/http"

	"git.sfxdx.ru/crystalline/wi-fi-backend/services/cloudtrax"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/twilio"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
	"github.com/labstack/echo"
)

type CryptoRouter struct {
	worldbitService  worldbit.Worldbit
	cloudtraxService cloudtrax.Cloudtrax
	twilioService    twilio.Twilio
}

func NewCryptoRouter(
	worldbitService worldbit.Worldbit,
	cloudtraxService cloudtrax.Cloudtrax,
	twilioService twilio.Twilio,
) CryptoRouter {
	return CryptoRouter{
		worldbitService:  worldbitService,
		cloudtraxService: cloudtraxService,
		twilioService:    twilioService,
	}
}

func (router CryptoRouter) Register(group *echo.Group) {
	group.POST("/payment", router.requestPayment)
}

func (router CryptoRouter) requestPayment(context echo.Context) error {
	request := new(PaymentRequest)
	if err := context.Bind(request); err != nil {
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
		Amount:         request.Amount * rate,
		SenderCurrency: request.Currency,
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

		voucherCode, err := router.cloudtraxService.CreateVoucher(
			request.NetworkID,
			cloudtrax.Voucher{
				Duration:  request.Voucher.Duration,
				MaxUsers:  request.Voucher.MaxUsers,
				UpLimit:   request.Voucher.UpLimit,
				DownLimit: request.Voucher.DownLimit,
				PurgeDays: request.Voucher.PurgeDays,
			},
		)
		if err != nil {
			log.Printf("Unable to generate voucher to network id %s. Error: %s\n", request.NetworkID, err)
			return
		}

		if err := router.twilioService.SendVoucher(request.PhoneNumber, voucherCode); err != nil {
			log.Printf("Unable to send voucher to phone number %s. Error: %s\n", request.PhoneNumber, err)
			return
		}

		log.Printf("Successfully generated voucher %s for user with phone number %s and network id %s\n",
			voucherCode,
			request.PhoneNumber,
			request.NetworkID,
		)
	}()

	return context.JSON(http.StatusOK, PaymentResponse{
		Address: exchange.Address,
		Amount:  exchange.Amount,
	})
}
