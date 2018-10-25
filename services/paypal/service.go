package paypal

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"git.sfxdx.ru/crystalline/wi-fi-backend/http"
	"github.com/pkg/errors"
)

type service struct {
	store           database.SalesStore
	config          Config
	token           string
	tokenExpiration time.Time
}

func New(store database.SalesStore, config Config) PayPal {
	return &service{
		store:  store,
		config: config,
	}
}

func (service *service) CheckSale(saleID string) error {
	if time.Now().After(service.tokenExpiration) {
		if err := service.resetToken(); err != nil {
			return err
		}
	}

	existingSale, err := service.store.ByPayPalSaleID(saleID)
	if err == nil && existingSale.PayPalSaleID == saleID {
		return errors.New("sale id is already used")
	}

	_, responseBody, err := http.Get(
		service.config.Host+"/payments/sale/"+saleID,
		http.Headers{
			"authorization": service.token,
		},
	)
	if err != nil {
		return err
	}

	response := new(saleDetails)

	if err := json.Unmarshal(responseBody, response); err != nil {
		return err
	}

	if response.State == COMPLETED {
		return nil
	}

	return errors.New("Payment is not complete or is denied")
}

func (service service) PersistSale(saleID string, voucher string) error {
	return service.store.Create(&database.UsedSale{
		PayPalSaleID: saleID,
		Voucher:      voucher,
	})
}

func (service *service) resetToken() error {
	authString := base64.StdEncoding.EncodeToString(
		[]byte(service.config.ClientID + ":" + service.config.Secret),
	)
	_, responseBody, err := http.Post(
		service.config.Host+"/oauth2/token",
		http.Headers{
			"content-type":  "application/x-www-form-urlencoded",
			"authorization": "Basic " + authString,
		},
		[]byte("grant_type=client_credentials"),
	)
	if err != nil {
		return err
	}

	response := new(tokenResponse)
	if err := json.Unmarshal(responseBody, response); err != nil {
		return err
	}
	if response.ExpiresIn == 0 {
		return errors.New("Invalid token expiration time")
	}

	service.tokenExpiration = time.Now().Add(time.Second * time.Duration(response.ExpiresIn))
	service.token = response.TokenType + " " + response.AccessToken

	return nil
}
