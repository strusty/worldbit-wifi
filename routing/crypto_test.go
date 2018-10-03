package routing

import (
	"testing"
	"time"

	"git.sfxdx.ru/crystalline/wi-fi-backend/services/cloudtrax"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewCryptoRouter(t *testing.T) {
	worldbitService := &WorldbitServiceMock{}
	cloudtraxService := &CloudtraxServiceMock{}
	twilioService := &TwilioServiceMock{}

	testRouter := CryptoRouter{
		worldbitService:  worldbitService,
		cloudtraxService: cloudtraxService,
		twilioService:    twilioService,
	}

	router := NewCryptoRouter(
		worldbitService,
		cloudtraxService,
		twilioService,
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

	cloudtraxService := &CloudtraxServiceMock{
		CreateVoucherFn: func(networkID string, voucher cloudtrax.Voucher) (string, error) {
			return "", nil
		},
	}

	twilioService := &TwilioServiceMock{
		SendVoucherFn: func(phoneNumber string, voucher string) error {
			return nil
		},
	}

	router := CryptoRouter{
		worldbitService:  worldbitService,
		cloudtraxService: cloudtraxService,
		twilioService:    twilioService,
	}

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContext()))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.True(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.True(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.MonitorExchangeStatusFnInvoked = false
	cloudtraxService.CreateVoucherFnInvoked = false
	twilioService.SendVoucherFnInvoked = false
	twilioService.SendVoucherFn = func(phoneNumber string, voucher string) error {
		return errors.New("test_error")
	}

	t.Run("Send voucher error", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContext()))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.True(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.True(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.MonitorExchangeStatusFnInvoked = false
	cloudtraxService.CreateVoucherFnInvoked = false
	cloudtraxService.CreateVoucherFn = func(networkID string, voucher cloudtrax.Voucher) (string, error) {
		return "", errors.New("test_error")
	}
	twilioService.SendVoucherFnInvoked = false

	t.Run("Create voucher error", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContext()))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.True(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.MonitorExchangeStatusFnInvoked = false
	worldbitService.MonitorExchangeStatusFn = func(statusURL string) error {
		return errors.New("test_error")
	}
	cloudtraxService.CreateVoucherFnInvoked = false

	t.Run("Monitor exchange status error", func(t *testing.T) {
		assert.NoError(t, router.requestPayment(generateContext()))
		time.Sleep(time.Second)
		assert.True(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.CreateExchangeFn = func(request worldbit.CreateExchangeRequest) (*worldbit.CreateExchangeResult, error) {
		return nil, errors.New("test_error")
	}
	worldbitService.MonitorExchangeStatusFnInvoked = false

	t.Run("Create exchange error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContext()))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.GetExchangeRateFn = func() (float64, error) {
		return 0, errors.New("test_error")
	}

	t.Run("Get exchange rate error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContext()))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	worldbitService.CreateAccountFn = func() (*worldbit.CreateAccountResponseData, error) {
		return nil, errors.New("test_error")
	}

	t.Run("Create account error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContext()))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

	t.Run("Invalid body error", func(t *testing.T) {
		assert.Error(t, router.requestPayment(generateContextWithInvalidBody()))
		assert.False(t, worldbitService.MonitorExchangeStatusFnInvoked)
		assert.False(t, cloudtraxService.CreateVoucherFnInvoked)
		assert.False(t, twilioService.SendVoucherFnInvoked)
	})

}
