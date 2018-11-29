package captcha

import (
	"encoding/json"
	"net/url"

	"github.com/strusty/worldbit-wifi/http"
)

type service struct {
	secret string
}

func New(secret string) Captcha {
	return service{
		secret: secret,
	}
}

func (controller service) CheckCaptcha(responseToken string) (bool, error) {
	type captchaResponse struct {
		Success bool `json:"success"`
	}

	form := url.Values{}
	form.Add("secret", controller.secret)
	form.Add("response", responseToken)

	_, responseData, err := http.Post(
		"https://www.google.com/recaptcha/api/siteverify",
		http.Headers{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		[]byte(form.Encode()),
	)

	if err != nil {
		return false, err
	}

	response := new(captchaResponse)

	if err := json.Unmarshal(responseData, response); err != nil {
		return false, err
	}

	return response.Success, nil
}
