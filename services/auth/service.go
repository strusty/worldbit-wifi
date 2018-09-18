package auth

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"git.sfxdx.ru/crystalline/wi-fi-backend/random"
	"time"
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

func (service service) CreateCode(request SendCodeRequest) (string, error) {
	confirmationCode := random.String(16)

	if err := service.store.Create(&database.Authentication{
		PhoneNumber:      request.PhoneNumber,
		ConfirmationCode: confirmationCode,
		ExpiryDate:       time.Now().Add(time.Second * time.Duration(service.expiration)),
	}); err != nil {
		return "", err
	}

	return confirmationCode, nil
}
