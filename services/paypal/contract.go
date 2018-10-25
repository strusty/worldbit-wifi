package paypal

type PayPal interface {
	CheckSale(saleID string) error
	PersistSale(saleID string, voucher string) error
}

type Config struct {
	Host     string
	ClientID string
	Secret   string
}
