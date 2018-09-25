package auth

type Auth interface {
	CreateCode(phoneNumber string) (string, error)
	VerifyCode(confirmationCode string) error
}
