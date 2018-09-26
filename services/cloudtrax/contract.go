package cloudtrax

type Cloudtrax interface {
	CreateVoucher(networkID string, voucher Voucher) (string, error)
}
