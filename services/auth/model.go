package auth

type SendCodeRequest struct {
	PhoneNumber string
}

type VerifyCodeRequest struct {
	ConfirmationCode string
}
