package auth

import (
	"time"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"git.sfxdx.ru/crystalline/wi-fi-backend/random"
	"github.com/pkg/errors"
)

type service struct {
	store      database.AuthenticationsStore
	expiration int64
}

func New(store database.AuthenticationsStore, expiration int64) Auth {
	return service{
		store:      store,
		expiration: expiration,
	}
}

func (service service) CreateCode(phoneNumber string) (string, error) {
	confirmationCode := random.String(8)

	if err := service.store.Create(&database.Authentication{
		PhoneNumber:      phoneNumber,
		ConfirmationCode: confirmationCode,
		ExpiryDate:       time.Now().Add(time.Second * time.Duration(service.expiration)),
	}); err != nil {
		return "", err
	}

	return confirmationCode, nil
}

func (service service) VerifyCode(confirmationCode string) error {
	authentication, err := service.store.ByConfirmationCode(confirmationCode)
	if err != nil {
		return err
	}

	if authentication.ExpiryDate.Before(time.Now()) {
		return errors.New("You code has already expired. Please request a new one.")
	}

	return nil
}
