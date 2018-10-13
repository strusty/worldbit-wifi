package routing

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/admins"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/radius"
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
	RadiusServiceMock
*/
type RadiusServiceMock struct {
	CreateCredentialsFnInvoked bool
	CreateCredentialsFn        func(plan radius.PricingPlan) (string, error)
}

func (mock *RadiusServiceMock) CreateCredentials(plan radius.PricingPlan) (string, error) {
	mock.CreateCredentialsFnInvoked = true
	return mock.CreateCredentialsFn(plan)
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

/*
	PricingPlanServiceMock
*/
type PricingPlanServiceMock struct {
	CreateFn func(plan *pricing_plans.PricingPlan) error
	UpdateFn func(plan *pricing_plans.PricingPlan) error
	DeleteFn func(id string) error
	AllFn    func() ([]pricing_plans.PricingPlan, error)
	ByIDFn   func(id string) (*pricing_plans.PricingPlan, error)
}

func (mock PricingPlanServiceMock) Create(plan *pricing_plans.PricingPlan) error {
	return mock.CreateFn(plan)
}

func (mock PricingPlanServiceMock) Update(plan *pricing_plans.PricingPlan) error {
	return mock.UpdateFn(plan)
}

func (mock PricingPlanServiceMock) Delete(id string) error {
	return mock.DeleteFn(id)
}

func (mock PricingPlanServiceMock) All() ([]pricing_plans.PricingPlan, error) {
	return mock.AllFn()
}

func (mock PricingPlanServiceMock) ByID(id string) (*pricing_plans.PricingPlan, error) {
	return mock.ByIDFn(id)
}

/*
	AdminServiceMock
*/
type AdminServiceMock struct {
	LoginFn          func(request admins.LoginRequest) (*admins.JWTResponse, error)
	ChangePasswordFn func(adminID string, request admins.ChangePasswordRequest) error
}

func (mock AdminServiceMock) Login(request admins.LoginRequest) (*admins.JWTResponse, error) {
	return mock.LoginFn(request)
}

func (mock AdminServiceMock) ChangePassword(adminID string, request admins.ChangePasswordRequest) error {
	return mock.ChangePasswordFn(adminID, request)
}
