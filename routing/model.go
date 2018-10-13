package routing

type SendCodeRequest struct {
	PhoneNumber string
}

type VerifyCodeRequest struct {
	ConfirmationCode string
	Captcha          string
}

type PaymentRequest struct {
	PhoneNumber   string `json:"phoneNumber"`
	NetworkID     string `json:"networkID"`
	Currency      string `json:"currency"`
	PricingPlanID string `json:"pricingPlanID"`
}

type PaymentResponse struct {
	Address string `json:"address"`
	Amount  string `json:"response"`
}
