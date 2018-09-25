package captcha

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestMakeCaptchaService(t *testing.T) {
	testService := service{
		secret: "secret",
	}
	service := New("secret")

	assert.Equal(t, testService, service)
}

func TestCaptcha_CheckCaptcha(t *testing.T) {
	service := New("secret")
	testResponseSuccessToken := "success_token"
	testResponseFailureToken := "failure_token"
	testResponseInvalidToken := "invalid_token"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://www.google.com/recaptcha/api/siteverify",
		func(request *http.Request) (*http.Response, error) {
			assert.Equal(t, "secret", request.FormValue("secret"))

			if request.FormValue("response") == testResponseSuccessToken {
				return httpmock.NewJsonResponse(http.StatusOK, map[string]bool{
					"success": true,
				})
			}
			if request.FormValue("response") == testResponseFailureToken {
				return httpmock.NewJsonResponse(http.StatusOK, map[string]bool{
					"success": false,
				})
			}
			if request.FormValue("response") == testResponseInvalidToken {
				return httpmock.NewStringResponse(http.StatusOK, "{"), nil
			}
			return nil, errors.New("request failed")
		},
	)

	passed, err := service.CheckCaptcha(testResponseSuccessToken)

	if assert.NoError(t, err) {
		assert.True(t, passed)
	}

	passed, err = service.CheckCaptcha(testResponseFailureToken)

	if assert.NoError(t, err) {
		assert.False(t, passed)
	}

	passed, err = service.CheckCaptcha(testResponseInvalidToken)

	if assert.Error(t, err) {
		assert.False(t, passed)
	}

	passed, err = service.CheckCaptcha("")

	if assert.Error(t, err) {
		assert.False(t, passed)
	}
}
