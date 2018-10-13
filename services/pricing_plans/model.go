package pricing_plans

type PricingPlan struct {
	ID        string  `json:"id"`
	AmountUSD float64 `json:"amountUSD"`
	Duration  int64   `json:"duration"`
	MaxUsers  int64   `json:"maxUsers"`
	UpLimit   int64   `json:"upLimit" `
	DownLimit int64   `json:"downLimit"`
	PurgeDays int64   `json:"purgeDays"`
}

type PricingPlanValidator struct {
	ID        string   `json:"id"`
	AmountUSD *float64 `json:"amountUSD" validate:"required"`
	Duration  *int64   `json:"duration" validate:"required"`
	MaxUsers  *int64   `json:"maxUsers" validate:"required"`
	UpLimit   *int64   `json:"upLimit" validate:"required"`
	DownLimit *int64   `json:"downLimit" validate:"required"`
	PurgeDays *int64   `json:"purgeDays" validate:"required"`
}

func (pricingPlanValidator *PricingPlanValidator) GetPricingPlan() (*PricingPlan, error) {
	return &PricingPlan{
		ID:        pricingPlanValidator.ID,
		AmountUSD: *pricingPlanValidator.AmountUSD,
		Duration:  *pricingPlanValidator.Duration,
		MaxUsers:  *pricingPlanValidator.MaxUsers,
		UpLimit:   *pricingPlanValidator.UpLimit,
		DownLimit: *pricingPlanValidator.DownLimit,
		PurgeDays: *pricingPlanValidator.PurgeDays,
	}, nil
}
