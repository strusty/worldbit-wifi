package routing

type SendCodeRequest struct {
	PhoneNumber string
}

type VerifyCodeRequest struct {
	ConfirmationCode string
	Captcha          string
}
