package cleaner

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/radius_database"
)

type AccountingStoreMock struct {
	SessionTimeSumFnInvoked             bool
	SessionTimeSumFn                    func(username string) (int64, error)
	DeleteAccountingByUsernameFnInvoked bool
	DeleteAccountingByUsernameFn        func(username string) error
}

func (mock *AccountingStoreMock) SessionTimeSum(username string) (int64, error) {
	mock.SessionTimeSumFnInvoked = true
	return mock.SessionTimeSumFn(username)
}

func (mock *AccountingStoreMock) DeleteAccountingByUsername(username string) error {
	mock.DeleteAccountingByUsernameFnInvoked = true
	return mock.DeleteAccountingByUsernameFn(username)
}

type CheckStoreMock struct {
	SessionChecksFnInvoked          bool
	SessionChecksFn                 func() ([]radius_database.Check, error)
	DeleteChecksByUsernameFnInvoked bool
	DeleteChecksByUsernameFn        func(username string) error
}

func (CheckStoreMock) Create(check *radius_database.Check) error {
	return errors.New("not_mocked")
}

func (mock *CheckStoreMock) SessionChecks() ([]radius_database.Check, error) {
	mock.SessionChecksFnInvoked = true
	return mock.SessionChecksFn()
}

func (mock *CheckStoreMock) DeleteChecksByUsername(username string) error {
	mock.DeleteChecksByUsernameFnInvoked = true
	return mock.DeleteChecksByUsernameFn(username)
}

type ReplyStoreMock struct {
	DeleteRepliesByUsernameFnInvoked bool
	DeleteRepliesByUsernameFn        func(username string) error
}

func (ReplyStoreMock) Create(check *radius_database.Reply) error {
	return errors.New("test_error")
}

func (mock *ReplyStoreMock) DeleteRepliesByUsername(username string) error {
	mock.DeleteRepliesByUsernameFnInvoked = true
	return mock.DeleteRepliesByUsernameFn(username)
}

func TestCleaner(t *testing.T) {
	accountingStoreMock := &AccountingStoreMock{
		SessionTimeSumFn: func(username string) (int64, error) {
			return 130, nil
		},
		DeleteAccountingByUsernameFn: func(username string) error {
			return nil
		},
	}

	checkStoreMock := &CheckStoreMock{
		SessionChecksFn: func() ([]radius_database.Check, error) {
			return []radius_database.Check{
				{
					Value: "119",
				},
				{
					Value: "119a",
				},
			}, nil
		},
		DeleteChecksByUsernameFn: func(username string) error {
			return nil
		},
	}

	replyStoreMock := &ReplyStoreMock{
		DeleteRepliesByUsernameFn: func(username string) error {
			return nil
		},
	}

	testService := service{
		accountingStore: accountingStoreMock,
		checkStore:      checkStoreMock,
		replyStore:      replyStoreMock,
	}

	service := New(
		accountingStoreMock,
		checkStoreMock,
		replyStoreMock,
	)
	t.Run("Initialization", func(t *testing.T) {
		if assert.Equal(t, testService, service) {
			go service.Start(time.Second * 2)
			time.Sleep(time.Second*2 + time.Millisecond*10)

			t.Run("Succesful cleanup", func(t *testing.T) {
				assert.True(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.True(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.True(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.True(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.True(t, checkStoreMock.SessionChecksFnInvoked)
			})

			accountingStoreMock.DeleteAccountingByUsernameFnInvoked = false
			replyStoreMock.DeleteRepliesByUsernameFnInvoked = false
			checkStoreMock.DeleteChecksByUsernameFnInvoked = false
			accountingStoreMock.SessionTimeSumFnInvoked = false
			checkStoreMock.SessionChecksFnInvoked = false

			t.Run("Reset state", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.False(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.False(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.False(t, checkStoreMock.SessionChecksFnInvoked)
			})

			accountingStoreMock.DeleteAccountingByUsernameFn = func(username string) error {
				return errors.New("test_error")
			}
			replyStoreMock.DeleteRepliesByUsernameFn = func(username string) error {
				return errors.New("test_error")
			}

			time.Sleep(time.Second * 2)
			t.Run("Accounting and replies deletion error", func(t *testing.T) {
				assert.True(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.True(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.True(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.True(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.True(t, checkStoreMock.SessionChecksFnInvoked)
			})

			accountingStoreMock.DeleteAccountingByUsernameFnInvoked = false
			replyStoreMock.DeleteRepliesByUsernameFnInvoked = false
			checkStoreMock.DeleteChecksByUsernameFnInvoked = false
			accountingStoreMock.SessionTimeSumFnInvoked = false
			checkStoreMock.SessionChecksFnInvoked = false

			t.Run("Reset state", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.False(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.False(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.False(t, checkStoreMock.SessionChecksFnInvoked)
			})

			checkStoreMock.DeleteChecksByUsernameFn = func(username string) error {
				return errors.New("test_error")
			}

			time.Sleep(time.Second * 2)

			t.Run("Checks deletion error", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.True(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.True(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.True(t, checkStoreMock.SessionChecksFnInvoked)
			})

			checkStoreMock.DeleteChecksByUsernameFnInvoked = false
			accountingStoreMock.SessionTimeSumFnInvoked = false
			checkStoreMock.SessionChecksFnInvoked = false

			t.Run("Reset state", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.False(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.False(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.False(t, checkStoreMock.SessionChecksFnInvoked)
			})

			accountingStoreMock.SessionTimeSumFn = func(username string) (int64, error) {
				return 0, errors.New("test_error")
			}

			time.Sleep(time.Second * 2)

			t.Run("Get session time sum error", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.False(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.True(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.True(t, checkStoreMock.SessionChecksFnInvoked)
			})

			accountingStoreMock.SessionTimeSumFnInvoked = false
			checkStoreMock.SessionChecksFnInvoked = false

			t.Run("Reset state", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.False(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.False(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.False(t, checkStoreMock.SessionChecksFnInvoked)
			})

			checkStoreMock.SessionChecksFn = func() ([]radius_database.Check, error) {
				return nil, errors.New("test_error")
			}

			time.Sleep(time.Second * 2)

			t.Run("Get session checks error", func(t *testing.T) {
				assert.False(t, accountingStoreMock.DeleteAccountingByUsernameFnInvoked)
				assert.False(t, replyStoreMock.DeleteRepliesByUsernameFnInvoked)
				assert.False(t, checkStoreMock.DeleteChecksByUsernameFnInvoked)
				assert.False(t, accountingStoreMock.SessionTimeSumFnInvoked)
				assert.True(t, checkStoreMock.SessionChecksFnInvoked)
			})
		}
	})

}
