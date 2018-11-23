package routing

type SendCodeRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

type VerifyCodeRequest struct {
	ConfirmationCode string `json:"confirmationCode"`
	Captcha          string `json:"captcha"`
}

type PaymentRequest struct {
	PhoneNumber   string `json:"phoneNumber"`
	NetworkID     string `json:"networkID"`
	Currency      string `json:"currency"`
	PricingPlanID string `json:"pricingPlanID"`
}

type PaymentResponse struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type PayPalVoucherRequest struct {
	SaleID        string `json:"saleID"`
	PricingPlanID string `json:"pricingPlanID"`
	PhoneNumber   string `json:"phoneNumber"`
}
