package paypal

import (
	"net/http"
	"testing"
	"time"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

type SalesStoreMock struct {
	CreateFn         func(sale *database.UsedSale) error
	ByPayPalSaleIDFn func(saleID string) (*database.UsedSale, error)
}

func (mock SalesStoreMock) Create(sale *database.UsedSale) error {
	return mock.CreateFn(sale)
}

func (mock SalesStoreMock) ByPayPalSaleID(saleID string) (*database.UsedSale, error) {
	return mock.ByPayPalSaleIDFn(saleID)
}

func TestNew(t *testing.T) {
	storeMock := &SalesStoreMock{}
	config := Config{
		Host:     "host",
		ClientID: "clientID",
		Secret:   "secret",
	}

	testService := &service{
		store:  storeMock,
		config: config,
	}

	payPalService := New(storeMock, config)

	assert.Equal(t, testService, payPalService)
}

func Test_service_CheckSale(t *testing.T) {
	storeMock := &SalesStoreMock{
		ByPayPalSaleIDFn: func(saleID string) (*database.UsedSale, error) {
			return nil, errors.New("test_error")
		},
	}

	service := New(storeMock, Config{
		Host: "http://test.com",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Success", func(t *testing.T) {
		httpmock.RegisterResponder(
			http.MethodGet,
			"http://test.com/payments/sale/test_id",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusOK, `{"state":"completed"}`), nil
			},
		)
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, &tokenResponse{
					AccessToken: "token",
					TokenType:   "Bearer",
					ExpiresIn:   5,
				})
			},
		)

		assert.NoError(t, service.CheckSale("test_id"))
	})

	t.Run("Incomplete transaction error", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, &tokenResponse{
					AccessToken: "token",
					TokenType:   "Bearer",
					ExpiresIn:   5,
				})
			},
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			"http://test.com/payments/sale/test_id",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusOK, `{"state":""}`), nil
			},
		)

		assert.Error(t, service.CheckSale("test_id"))
	})

	t.Run("Invalid json error. Sale endpoint.", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, &tokenResponse{
					AccessToken: "token",
					TokenType:   "Bearer",
					ExpiresIn:   5,
				})
			},
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			"http://test.com/payments/sale/test_id",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusOK, `{"state":""`), nil
			},
		)

		assert.Error(t, service.CheckSale("test_id"))
	})

	t.Run("Sale endpoint error.", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, &tokenResponse{
					AccessToken: "token",
					TokenType:   "Bearer",
					ExpiresIn:   5,
				})
			},
		)
		httpmock.RegisterResponder(
			http.MethodGet,
			"http://test.com/payments/sale/test_id",
			func(request *http.Request) (*http.Response, error) {
				return nil, errors.New("test_error")
			},
		)

		assert.Error(t, service.CheckSale("test_id"))
	})

	storeMock.ByPayPalSaleIDFn = func(saleID string) (*database.UsedSale, error) {
		return &database.UsedSale{
			PayPalSaleID: saleID,
		}, nil
	}

	t.Run("The sale is already used error", func(t *testing.T) {
		assert.Error(t, service.CheckSale("test_id"))
	})

	time.Sleep(time.Second * 6)

	t.Run("Invalid token expiration date", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewJsonResponse(http.StatusOK, &tokenResponse{
					AccessToken: "token",
					TokenType:   "Bearer",
					ExpiresIn:   0,
				})
			},
		)
		assert.Error(t, service.CheckSale("test_id"))
	})

	t.Run("Invalid json. Token endpoint", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusOK, "{"), nil
			},
		)
		assert.Error(t, service.CheckSale("test_id"))
	})

	t.Run("Token endpoint error", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			http.MethodPost,
			"http://test.com/oauth2/token",
			func(request *http.Request) (*http.Response, error) {
				return nil, errors.New("test_error")
			},
		)
		assert.Error(t, service.CheckSale("test_id"))
	})
}

func Test_service_PersistSale(t *testing.T) {
	storeMock := &SalesStoreMock{
		CreateFn: func(sale *database.UsedSale) error {
			return nil
		},
	}

	service := New(storeMock, Config{
		Host: "http://test.com",
	})

	assert.NoError(t, service.PersistSale("id", "voucher"))

	storeMock.CreateFn = func(sale *database.UsedSale) error {
		return errors.New("test_error")
	}

	assert.Error(t, service.PersistSale("id", "voucher"))
}
