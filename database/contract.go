package database

type AuthenticationsStore interface {
	Create(entity *Authentication) error
	ByPhoneNumber(phoneNumber string) (*Authentication, error)
	ByConfirmationCode(confirmationCode string) (*Authentication, error)
}
