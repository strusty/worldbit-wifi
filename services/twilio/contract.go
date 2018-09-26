package twilio

type Twilio interface {
	SendConfirmationCode(phoneNumber string, confirmationCode string) error
	SendVoucher(phoneNumber string, voucher string) error
}
