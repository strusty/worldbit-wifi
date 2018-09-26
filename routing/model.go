package routing

type SendCodeRequest struct {
	PhoneNumber string
}

type VerifyCodeRequest struct {
	ConfirmationCode string
	Captcha          string
}

type PaymentRequest struct {
	PhoneNumber string  `json:"phoneNumber"`
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
	NetworkID   string  `json:"networkID"`
	Voucher     Voucher `json:"voucher"`
}

type Voucher struct {
	Duration  int64 `json:"duration"`
	MaxUsers  int64 `json:"maxUsers"`
	UpLimit   int64 `json:"upLimit"`
	DownLimit int64 `json:"downLimit"`
	PurgeDays int64 `json:"purgeDays"`
}

type PaymentResponse struct {
	Address string `json:"address"`
	Amount  string `json:"response"`
}
