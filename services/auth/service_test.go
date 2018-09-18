package auth

import (
	"testing"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"time"
)

type AuthenticationsStoreMock struct {
	CreateFn             func(entity *database.Authentication) error
	ByPhoneNumberFn      func(phoneNumber string) (*database.Authentication, error)
	ByConfirmationCodeFn func(confirmationCode string) (*database.Authentication, error)
}

func (mock AuthenticationsStoreMock) Create(entity *database.Authentication) error {
	return mock.CreateFn(entity)
}

func (mock AuthenticationsStoreMock) ByPhoneNumber(phoneNumber string) (*database.Authentication, error) {
	return mock.ByPhoneNumberFn(phoneNumber)
}

func (mock AuthenticationsStoreMock) ByConfirmationCode(confirmationCode string) (*database.Authentication, error) {
	return mock.ByConfirmationCodeFn(confirmationCode)
}

func TestNew(t *testing.T) {
	store := AuthenticationsStoreMock{}

	testService := service{
		store:      &store,
		expiration: 10,
	}

	service := New(&store, 10)

	assert.Equal(t, testService, service)
}

func Test_service_CreateCode(t *testing.T) {
	service := service{
		store: AuthenticationsStoreMock{
			CreateFn: func(entity *database.Authentication) error {
				return nil
			},
		},
	}

	code, err := service.CreateCode(SendCodeRequest{})
	if assert.NoError(t, err) {
		newCode, err := service.CreateCode(SendCodeRequest{})
		if assert.NoError(t, err) {
			assert.NotEqual(t, code, newCode)
		}
	}

	service.store = AuthenticationsStoreMock{
		CreateFn: func(entity *database.Authentication) error {
			return errors.New("test_error")
		},
	}

	_, err = service.CreateCode(SendCodeRequest{})
	assert.Error(t, err)
}

func Test_service_VerifyCode(t *testing.T) {
	service := service{
		store: AuthenticationsStoreMock{
			ByConfirmationCodeFn: func(confirmationCode string) (*database.Authentication, error) {
				return &database.Authentication{
					ExpiryDate: time.Now().Add(time.Hour),
				}, nil
			},
		},
	}

	assert.NoError(t, service.VerifyCode(VerifyCodeRequest{}))

	service.store = AuthenticationsStoreMock{
		ByConfirmationCodeFn: func(confirmationCode string) (*database.Authentication, error) {
			return &database.Authentication{}, nil
		},
	}

	assert.Error(t, service.VerifyCode(VerifyCodeRequest{}))

	service.store = AuthenticationsStoreMock{
		ByConfirmationCodeFn: func(confirmationCode string) (*database.Authentication, error) {
			return nil, errors.New("test_error")
		},
	}

	assert.Error(t, service.VerifyCode(VerifyCodeRequest{}))
}
