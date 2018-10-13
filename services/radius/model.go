package radius

type PricingPlan struct {
	Duration  int64
	MaxUsers  int64
	UpLimit   int64
	DownLimit int64
	PurgeDays int64
}
