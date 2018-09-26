package cloudtrax

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	testService := service{
		APIKey: "test_key",
		Secret: "test_secret",
		Host:   "https://test.com",
	}

	service := New(
		"test_key",
		"test_secret",
		"https://test.com",
	)

	assert.Equal(t, testService, service)
}

func Test_service_CreateVoucher(t *testing.T) {
	testService := service{
		APIKey: "test_key",
		Secret: "test_secret",
		Host:   "https://test.com",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test.com/voucher/network/test_id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{ "vouchers": [{ "voucher_code": "test_code" }]}`,
			), nil
		},
	)

	code, err := testService.CreateVoucher("test_id", Voucher{})

	assert.Equal(t, "test_code", code)
	assert.NoError(t, err)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test.com/voucher/network/test_id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{ }`,
			), nil
		},
	)
	code, err = testService.CreateVoucher("test_id", Voucher{})

	assert.Empty(t, code)
	if assert.Error(t, err) {
		assert.Equal(t, "Unknown error happened", err.Error())
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test.com/voucher/network/test_id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{ "errors": [{"code": 10, "message": "test_message"}, {"code": 10023, "message": "other_test_message"}] }`,
			), nil
		},
	)
	code, err = testService.CreateVoucher("test_id", Voucher{})

	assert.Empty(t, code)
	if assert.Error(t, err) {
		assert.Equal(t, "code=10 message=test_message;code=10023 message=other_test_message;", err.Error())
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test.com/voucher/network/test_id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{ "errors": [] }`,
			), nil
		},
	)
	code, err = testService.CreateVoucher("test_id", Voucher{})

	assert.Empty(t, code)
	if assert.Error(t, err) {
		assert.Equal(t, "", err.Error())
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test.com/voucher/network/test_id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{ "errors": [] `,
			), nil
		},
	)
	code, err = testService.CreateVoucher("test_id", Voucher{})

	assert.Empty(t, code)
	assert.Error(t, err)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test.com/voucher/network/test_id",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)
	code, err = testService.CreateVoucher("test_id", Voucher{})

	assert.Empty(t, code)
	assert.Error(t, err)
}
