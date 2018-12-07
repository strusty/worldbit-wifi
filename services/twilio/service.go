package twilio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/strusty/worldbit-wifi/http"
)

type service struct {
	accessToken                     string
	endpoint                        string
	from                            string
	confirmationCodeMessageTemplate string
	voucherMessageTemplate          string
}

func New(
	host string,
	sid string,
	token string,
	from string,
	confirmationCodeMessageTemplate string,
	voucherMessageTemplate string,
) Twilio {
	return service{
		accessToken: base64.StdEncoding.EncodeToString(
			[]byte(sid + ":" + token),
		),
		endpoint:                        host + "/Accounts/" + sid + "/Messages.json",
		from:                            from,
		confirmationCodeMessageTemplate: confirmationCodeMessageTemplate,
		voucherMessageTemplate:          voucherMessageTemplate,
	}
}

func (service service) SendConfirmationCode(phoneNumber string, confirmationCode string) error {
	return service.sendTemplatedMessage(phoneNumber, service.confirmationCodeMessageTemplate, confirmationCode)
}

func (service service) SendVoucher(phoneNumber string, voucher string) error {
	return service.sendTemplatedMessage(phoneNumber, service.voucherMessageTemplate, voucher)
}

func (service service) sendTemplatedMessage(phoneNumber string, template string, values ...interface{}) error {
	type Response struct {
		ErrorCode    *int64  `json:"code"`
		ErrorMessage *string `json:"message"`
	}

	form := url.Values{}
	form.Add("Body", fmt.Sprintf(template, values...))
	form.Add("To", "+"+phoneNumber)
	form.Add("From", service.from)

	_, responseData, err := http.Post(
		service.endpoint,
		http.Headers{
			"Authorization": "Basic " + service.accessToken,
			"Content-Type":  "application/x-www-form-urlencoded",
		},
		[]byte(form.Encode()),
	)
	if err != nil {
		return err
	}

	response := new(Response)

	if err := json.Unmarshal(responseData, response); err != nil {
		return err
	}

	if response.ErrorMessage != nil {
		return errors.New(*response.ErrorMessage)
	} else if response.ErrorCode != nil {
		return errors.Errorf("Error code: %d", *response.ErrorCode)
	}

	return nil
}
