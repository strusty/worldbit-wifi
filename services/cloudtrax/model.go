package cloudtrax

type CreateVouchersRequest struct {
	DesiredVouchers []Voucher `json:"desired_vouchers"`
}

type CreateVouchersResponse struct {
	Vouchers []Voucher `json:"vouchers"`
	Errors   []Error   `json:"errors"`
}

type Voucher struct {
	Code        string `json:"code"`
	Duration    int64  `json:"duration"`
	MaxUsers    int64  `json:"max_users,omitempty"`
	UpLimit     int64  `json:"up_limit,omitempty"`
	DownLimit   int64  `json:"down_limit,omitempty"`
	PurgeDays   int64  `json:"purge_days"`
	VoucherCode string `json:"voucher_code,omitempty"`
}

type Error struct {
	Code    int64
	Message string
}
