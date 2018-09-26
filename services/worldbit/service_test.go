package worldbit

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	config := Config{
		APIKey:            "APIKey",
		APISecret:         "APISecret",
		MerchantID:        "MerchantID",
		Host:              "Host",
		MonitoringTimeout: 10,
	}

	testService := service{
		Config: config,
	}

	service := New(config)

	assert.Equal(t, testService, service)
}

func Test_service_CreateExchange(t *testing.T) {
	service := service{
		Config{
			Host: "http://test.com",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_exchange",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
"status": true,
"result": {
    "amount": "0.00387000",
    "address": "0x2593BBc4001E50d927277e92bdc93f40E2b7C70f",
    "txn_id": "aPQbl8gg3jLpF6Uh7ZenbIaHUXBlyEAcq8o0sr6UyE3XrNbN",
    "status_url": "https://api.worldbit.com/status/aPQbl8gg3jLpF6Uh7ZenbIaHUXBlyEAcq8o0sr6UyE3XrNbN/0x2593BBc4001E50d927277e92bdc93f40E2b7C70f"
}}`,
			), nil
		},
	)

	t.Run("Success", func(t *testing.T) {
		result, err := service.CreateExchange(CreateExchangeRequest{})
		if assert.NoError(t, err) && assert.NotNil(t, result) {
			assert.Equal(t, "0.00387000", result.Amount)
			assert.Equal(t, "0x2593BBc4001E50d927277e92bdc93f40E2b7C70f", result.Address)
			assert.Equal(t, "aPQbl8gg3jLpF6Uh7ZenbIaHUXBlyEAcq8o0sr6UyE3XrNbN", result.TransactionID)
			assert.Equal(t, "https://api.worldbit.com/status/aPQbl8gg3jLpF6Uh7ZenbIaHUXBlyEAcq8o0sr6UyE3XrNbN/0x2593BBc4001E50d927277e92bdc93f40E2b7C70f", result.StatusURL)
		}
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_exchange",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": false}`,
			), nil
		},
	)

	t.Run("Unknown error", func(t *testing.T) {
		result, err := service.CreateExchange(CreateExchangeRequest{})
		if assert.Error(t, err) && assert.Nil(t, result) {
			assert.Equal(t, "Unknown error occurred", err.Error())
		}
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_exchange",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": false, "msg": "test_error"}`,
			), nil
		},
	)

	t.Run("Error with message", func(t *testing.T) {
		result, err := service.CreateExchange(CreateExchangeRequest{})
		if assert.Error(t, err) && assert.Nil(t, result) {
			assert.Equal(t, "test_error", err.Error())
		}
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_exchange",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": false, "msg": "test_error"`,
			), nil
		},
	)

	t.Run("Invalid response body", func(t *testing.T) {
		result, err := service.CreateExchange(CreateExchangeRequest{})
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_exchange",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	t.Run("Request failed", func(t *testing.T) {
		result, err := service.CreateExchange(CreateExchangeRequest{})
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func Test_service_CreateAccount(t *testing.T) {
	service := service{
		Config{
			Host: "http://test.com",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_account",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
"status": true,
"data": {
        "symbol": "BTC",
        "address": "2MvRDW2zuaKY8yGgBfjFuwvVmaYvdHp7XVd",
        "accountname": "BTCjRtG2vDj70OdyTfLl018lHZAEWza7eEuBPV1orzm8",
        "buyer_email": "buyer@email.com"
    }}`,
			), nil
		},
	)

	t.Run("Success", func(t *testing.T) {
		result, err := service.CreateAccount(CreateAccountRequest{})
		if assert.NoError(t, err) && assert.NotNil(t, result) {
			assert.Equal(t, "2MvRDW2zuaKY8yGgBfjFuwvVmaYvdHp7XVd", result.Address)
			assert.Equal(t, "buyer@email.com", result.BuyerEmail)
			assert.Equal(t, "BTC", result.Symbol)
			assert.Equal(t, "BTCjRtG2vDj70OdyTfLl018lHZAEWza7eEuBPV1orzm8", result.AccountName)
		}
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_account",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": false}`,
			), nil
		},
	)

	t.Run("Unknown error", func(t *testing.T) {
		result, err := service.CreateAccount(CreateAccountRequest{})
		if assert.Error(t, err) && assert.Nil(t, result) {
			assert.Equal(t, "Unknown error occurred", err.Error())
		}
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_account",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": false, "msg": "test_error"}`,
			), nil
		},
	)

	t.Run("Error with message", func(t *testing.T) {
		result, err := service.CreateAccount(CreateAccountRequest{})
		if assert.Error(t, err) && assert.Nil(t, result) {
			assert.Equal(t, "test_error", err.Error())
		}
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_account",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": false, "msg": "test_error"`,
			), nil
		},
	)

	t.Run("Invalid response body", func(t *testing.T) {
		result, err := service.CreateAccount(CreateAccountRequest{})
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com/create_account",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	t.Run("Request failed", func(t *testing.T) {
		result, err := service.CreateAccount(CreateAccountRequest{})
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func Test_service_GetExchangeRate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://intel.worldbit.com/wbtprices.php",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status":"1","ethbtc":"0.03425","ethusd":"221.61","wbteth":0.00303,"wbtusd":0.6714783000000001}`,
			), nil
		},
	)

	t.Run("Success", func(t *testing.T) {
		rate, err := service{}.GetExchangeRate()
		assert.NoError(t, err)
		assert.Equal(t, 0.6714783000000001, rate)
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://intel.worldbit.com/wbtprices.php",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status":"0"}`,
			), nil
		},
	)

	t.Run("Error", func(t *testing.T) {
		rate, err := service{}.GetExchangeRate()
		assert.Error(t, err)
		assert.Equal(t, 0.0, rate)
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://intel.worldbit.com/wbtprices.php",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status":"0"`,
			), nil
		},
	)

	t.Run("Invalid response body", func(t *testing.T) {
		rate, err := service{}.GetExchangeRate()
		assert.Error(t, err)
		assert.Equal(t, 0.0, rate)
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://intel.worldbit.com/wbtprices.php",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	t.Run("Request failed", func(t *testing.T) {
		rate, err := service{}.GetExchangeRate()
		assert.Error(t, err)
		assert.Equal(t, 0.0, rate)
	})

}

func Test_service_MonitorExchangeStatus(t *testing.T) {
	service := service{
		Config{
			MonitoringTimeout: 30,
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": true, "data": { "status": 1 }}`,
			), nil
		},
	)

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, service.MonitorExchangeStatus("https://test.com"))
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": true, "data": { "status": -1 }}`,
			), nil
		},
	)

	t.Run("Worldbit timeout", func(t *testing.T) {
		assert.Error(t, service.MonitorExchangeStatus("https://test.com"))

	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": true, "data": { "status": -100 }}`,
			), nil
		},
	)

	t.Run("No error timeout", func(t *testing.T) {
		assert.Error(t, service.MonitorExchangeStatus("https://test.com"))

	})

	service.MonitoringTimeout = 19
	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"status": true, "data": { "status": -100 }`,
			), nil
		},
	)

	t.Run("Invalid response body timeout", func(t *testing.T) {
		assert.Error(t, service.MonitorExchangeStatus("https://test.com"))
	})

	service.MonitoringTimeout = 19
	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test.com",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	t.Run("Request failed timeout", func(t *testing.T) {
		assert.Error(t, service.MonitorExchangeStatus("https://test.com"))
	})
}
