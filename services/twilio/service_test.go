package twilio

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	host := "http://test.com"
	sid := "sid"

	testService := service{
		accessToken: base64.StdEncoding.EncodeToString(
			[]byte(sid + ":" + "token"),
		),
		endpoint:                        host + "/Accounts/" + sid + "/Messages.json",
		from:                            "from",
		confirmationCodeMessageTemplate: "template",
		voucherMessageTemplate:          "other_template",
	}

	service := New(
		host,
		sid,
		"token",
		"from",
		"template",
		"other_template",
	)

	assert.Equal(t, testService, service)
}

func Test_service_SendConfirmationCode(t *testing.T) {
	service := service{
		endpoint: "http://test.com",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{}`,
			), nil
		},
	)
	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, service.SendConfirmationCode("", ""))
		assert.NoError(t, service.SendVoucher("", ""))
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"code": 303}`,
			), nil
		},
	)
	t.Run("Only code error", func(t *testing.T) {
		assert.Error(t, service.SendConfirmationCode("", ""))
		assert.Error(t, service.SendVoucher("", ""))
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"code": 303, "message": "test_error"}`,
			), nil
		},
	)
	t.Run("Code and message error", func(t *testing.T) {
		err := service.SendConfirmationCode("", "")
		assert.Error(t, err)
		assert.Equal(t, "test_error", err.Error())

		err = service.SendVoucher("", "")
		assert.Error(t, err)
		assert.Equal(t, "test_error", err.Error())
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"code": 3`,
			), nil
		},
	)

	t.Run("Invalid response body", func(t *testing.T) {
		assert.Error(t, service.SendConfirmationCode("", ""))
		assert.Error(t, service.SendVoucher("", ""))
	})

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://test.com",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)
	t.Run("Request failed", func(t *testing.T) {
		assert.Error(t, service.SendConfirmationCode("", ""))
		assert.Error(t, service.SendVoucher("", ""))
	})
}
