package cloudtrax

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"git.sfxdx.ru/crystalline/wi-fi-backend/http"
	"git.sfxdx.ru/crystalline/wi-fi-backend/random"
	"github.com/pkg/errors"
	"time"
)

type service struct {
	APIKey string
	Secret string
	Host   string
}

func New(apiKey string, secret string, host string) Cloudtrax {
	return service{
		APIKey: apiKey,
		Secret: secret,
		Host:   host,
	}
}

func (service service) CreateVoucher(networkID string, voucher Voucher) (string, error) {
	authorization := fmt.Sprintf("key=%s,timestamp=%d,nonce=%s",
		service.APIKey,
		time.Now().Unix(),
		random.String(8),
	)

	requestBytes, err := json.Marshal(CreateVouchersRequest{
		DesiredVouchers: []Voucher{voucher},
	})
	if err != nil {
		return "", err
	}

	endpoint := "/voucher/network/" + networkID

	signatureString := authorization + endpoint + string(requestBytes)

	hasher := hmac.New(sha256.New, []byte(service.Secret))
	hasher.Write([]byte(signatureString))
	hashedSignatureString := hasher.Sum(nil)

	_, responseBytes, err := http.Post(
		service.Host+endpoint,
		http.Headers{
			"Content-Type":         "application/json",
			"OpenMesh-API-Version": "1",
			"Authorization":        authorization,
			"Signature":            hex.EncodeToString(hashedSignatureString),
		},
		requestBytes)
	if err != nil {
		return "", err
	}

	response := new(CreateVouchersResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return "", err
	}

	if response.Errors != nil {
		errorMessage := ""

		for _, cloudtraxError := range response.Errors {
			errorMessage += fmt.Sprintf("code=%d message=%s;", cloudtraxError.Code, cloudtraxError.Message)
		}

		return "", errors.New(errorMessage)
	}

	if response.Vouchers != nil {
		return response.Vouchers[0].VoucherCode, nil
	}

	return "", errors.New("Unknown error happened")
}
