package worldbit

import (
	"encoding/json"
	"git.sfxdx.ru/crystalline/wi-fi-backend/http"
	"github.com/pkg/errors"
	"math"
	"time"
)

type service struct {
	Config
}

func New(config Config) Worldbit {
	return service{
		Config: config,
	}
}

func (service service) CreateExchange(request CreateExchangeRequest) (*CreateExchangeResult, error) {
	request.Merchant = service.Config.MerchantID
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, responseBytes, err := http.Post(service.Host+"/create_exchange", http.Headers{
		"Content-Type": "application/json",
	}, requestBytes)
	if err != nil {
		return nil, err
	}

	response := new(CreateExchangeResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	if !response.Status {
		if response.Message != nil {
			return nil, errors.New(*response.Message)
		} else {
			return nil, errors.New("Unknown error occurred")
		}
	}

	return &response.Result, nil
}

func (service service) CreateAccount(request CreateAccountRequest) (*CreateAccountResponseData, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, responseBytes, err := http.Post(service.Host+"/create_account", http.Headers{
		"Content-Type": "application/json",
	}, requestBytes)
	if err != nil {
		return nil, err
	}

	response := new(CreateAccountResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	if !response.Status {
		if response.Message != nil {
			return nil, errors.New(*response.Message)
		} else {
			return nil, errors.New("Unknown error occurred")
		}
	}

	return &response.Data, nil
}

func (service service) GetExchangeRate() (float64, error) {
	_, responseBytes, err := http.Get("https://intel.worldbit.com/wbtprices.php", http.Headers{})
	if err != nil {
		return 0, err
	}

	response := new(ExchangeRateResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return 0, err
	}

	if response.WbtUSD == 0 {
		return 0, errors.New("Unknown error")
	}

	return response.WbtUSD, nil
}

func (service service) MonitorExchangeStatus(statusURL string) (bool, error) {
	startingDate := time.Now()
	startingPower := 3
	for {
		if time.Now().Unix()-startingDate.Unix() >= int64(time.Second.Seconds())*service.MonitoringTimeout {
			return false, errors.New("Operation timed out")
		}
		time.Sleep(time.Duration(math.Exp(float64(startingPower))) * time.Second)
		startingPower += 1

		_, responseBytes, err := http.Get(statusURL, http.Headers{})
		if err != nil {
			continue
		}

		response := new(ExchangeStatus)
		if err := json.Unmarshal(responseBytes, response); err != nil {
			continue
		}

		switch response.Data.Status {
		case Timeout:
			return false, errors.New("Operation timed out")
		case Completed, PreInformCompleted, ConfirmedDeposit, ReceivedDeposit:
			return true, nil
		}
	}
}
