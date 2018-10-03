package routing

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/cloudtrax"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
)

/*
	AuthServiceMock
*/
type AuthServiceMock struct {
	CreateCodeFn func(phoneNumber string) (string, error)
	VerifyCodeFn func(confirmationCode string) error
}

func (mock AuthServiceMock) CreateCode(phoneNumber string) (string, error) {
	return mock.CreateCodeFn(phoneNumber)
}

func (mock AuthServiceMock) VerifyCode(confirmationCode string) error {
	return mock.VerifyCodeFn(confirmationCode)
}

/*
	TwilioServiceMock
*/
type TwilioServiceMock struct {
	SendConfirmationCodeFn func(phoneNumber string, confirmationCode string) error
	SendVoucherFnInvoked   bool
	SendVoucherFn          func(phoneNumber string, voucher string) error
}

func (mock TwilioServiceMock) SendConfirmationCode(phoneNumber string, confirmationCode string) error {
	return mock.SendConfirmationCodeFn(phoneNumber, confirmationCode)
}

func (mock *TwilioServiceMock) SendVoucher(phoneNumber string, voucher string) error {
	mock.SendVoucherFnInvoked = true
	return mock.SendVoucherFn(phoneNumber, voucher)
}

type CaptchaServiceMock struct {
	CheckCaptchaFn func(responseToken string) (bool, error)
}

func (mock CaptchaServiceMock) CheckCaptcha(responseToken string) (bool, error) {
	return mock.CheckCaptchaFn(responseToken)
}

/*
	CloudtraxServiceMock
*/
type CloudtraxServiceMock struct {
	CreateVoucherFnInvoked bool
	CreateVoucherFn        func(networkID string, voucher cloudtrax.Voucher) (string, error)
}

func (mock *CloudtraxServiceMock) CreateVoucher(networkID string, voucher cloudtrax.Voucher) (string, error) {
	mock.CreateVoucherFnInvoked = true
	return mock.CreateVoucherFn(networkID, voucher)
}

/*
	WorldbitServiceMock
*/
type WorldbitServiceMock struct {
	CreateExchangeFn               func(request worldbit.CreateExchangeRequest) (*worldbit.CreateExchangeResult, error)
	CreateAccountFn                func() (*worldbit.CreateAccountResponseData, error)
	GetExchangeRateFn              func() (float64, error)
	MonitorExchangeStatusFnInvoked bool
	MonitorExchangeStatusFn        func(statusURL string) error
}

func (mock WorldbitServiceMock) CreateExchange(request worldbit.CreateExchangeRequest) (*worldbit.CreateExchangeResult, error) {
	return mock.CreateExchangeFn(request)
}

func (mock WorldbitServiceMock) CreateAccount() (*worldbit.CreateAccountResponseData, error) {
	return mock.CreateAccountFn()
}

func (mock WorldbitServiceMock) GetExchangeRate() (float64, error) {
	return mock.GetExchangeRateFn()
}

func (mock *WorldbitServiceMock) MonitorExchangeStatus(statusURL string) error {
	mock.MonitorExchangeStatusFnInvoked = true
	return mock.MonitorExchangeStatusFn(statusURL)
}
