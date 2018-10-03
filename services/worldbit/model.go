package worldbit

type createExchangeRequest struct {
	Amount     float64 `json:"amount"`
	Currency1  string  `json:"currency1"`
	Currency2  string  `json:"currency2"`
	BuyerEmail string  `json:"buyer_email"`
	Address    string  `json:"address"`
	Merchant   string  `json:"merchant"`
}

type CreateExchangeRequest struct {
	Amount         float64
	SenderCurrency string
	Address        string
}

type CreateExchangeResponse struct {
	Status  bool                 `json:"status"`
	Result  CreateExchangeResult `json:"result"`
	Message *string              `json:"msg"`
}

type CreateExchangeResult struct {
	Amount        string `json:"amount"`
	Address       string `json:"address"`
	TransactionID string `json:"txn_id"`
	StatusURL     string `json:"status_url"`
}

type ExchangeStatus struct {
	Status bool               `json:"status"`
	Data   ExchangeStatusData `json:"data"`
}

type ExchangeStatusData struct {
	Status ExchangeStatusEnum `json:"status"`
}

type ExchangeStatusEnum int64

const (
	Timeout            = -1
	ReceivedDeposit    = 1
	ConfirmedDeposit   = 4
	PreInformCompleted = 99
	Completed          = 100
)

type createAccountRequest struct {
	Coin       string `json:"coin"`
	BuyerEmail string `json:"buyer_email"`
}

type CreateAccountResponse struct {
	Status  bool                      `json:"status"`
	Data    CreateAccountResponseData `json:"data"`
	Message *string                   `json:"msg"`
}

type CreateAccountResponseData struct {
	Symbol      string `json:"symbol"`
	Address     string `json:"address"`
	AccountName string `json:"accountname"`
	BuyerEmail  string `json:"buyer_email"`
}

type ExchangeRateResponse struct {
	WbtUSD float64 `json:"wbtusd"`
}
