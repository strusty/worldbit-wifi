package twilio

type Twilio interface {
	SendConfirmationCode(phoneNumber string, confirmationCode string) error
}
