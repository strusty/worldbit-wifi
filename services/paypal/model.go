package paypal

type saleDetails struct {
	ID    string    `json:"id"`
	State saleState `json:"state"`
}

type saleState string

const (
	COMPLETED saleState = "completed"
	PENDING   saleState = "pending"
	DENIED    saleState = "denied"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
