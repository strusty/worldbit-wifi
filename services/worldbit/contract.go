package worldbit

type Worldbit interface {
	CreateExchange(request CreateExchangeRequest) (*CreateExchangeResult, error)
	CreateAccount() (*CreateAccountResponseData, error)
	GetExchangeRate() (float64, error)
	MonitorExchangeStatus(statusURL string) error
}

type Config struct {
	APIKey            string
	APISecret         string
	MerchantID        string
	Host              string
	MonitoringTimeout int64
	DefaultCurrency   string
	DefaultEmail      string
}
