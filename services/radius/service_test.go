package radius

import (
	"testing"

	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type CheckStoreMock struct {
	CreateFn func(check *radius_database.Check) error
}

func (mock CheckStoreMock) Create(check *radius_database.Check) error {
	return mock.CreateFn(check)
}

type ReplyStoreMock struct {
	CreateFn func(check *radius_database.Reply) error
}

func (mock ReplyStoreMock) Create(check *radius_database.Reply) error {
	return mock.CreateFn(check)
}

func TestService(t *testing.T) {
	checkStore := &CheckStoreMock{
		CreateFn: func(check *radius_database.Check) error {
			return nil
		},
	}

	replyStore := &ReplyStoreMock{
		CreateFn: func(check *radius_database.Reply) error {
			return nil
		},
	}

	testService := service{
		checkStore: checkStore,
		replyStore: replyStore,
	}

	service := New(
		checkStore,
		replyStore,
	)

	t.Run("Initialization", func(t *testing.T) {
		if assert.Equal(t, testService, service) {
			t.Run("Success", func(t *testing.T) {
				code, err := service.CreateCredentials(PricingPlan{})
				assert.NoError(t, err)
				assert.NotEmpty(t, code)
			})

			replyStore.CreateFn = func(check *radius_database.Reply) error {
				return errors.New("test_error")
			}

			t.Run("Create reply error", func(t *testing.T) {
				_, err := service.CreateCredentials(PricingPlan{})
				assert.Error(t, err)
			})

			checkStore.CreateFn = func(check *radius_database.Check) error {
				return errors.New("test_error")
			}

			t.Run("Create check error", func(t *testing.T) {
				_, err := service.CreateCredentials(PricingPlan{})
				assert.Error(t, err)
			})
		}
	})

}
