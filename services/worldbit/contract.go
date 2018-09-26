package worldbit

type Worldbit interface {
	CreateExchange(request CreateExchangeRequest) (*CreateExchangeResult, error)
	CreateAccount(request CreateAccountRequest) (*CreateAccountResponseData, error)
	GetExchangeRate() (float64, error)
	MonitorExchangeStatus(statusURL string) error
}

type Config struct {
	APIKey            string
	APISecret         string
	MerchantID        string
	Host              string
	MonitoringTimeout int64
}
