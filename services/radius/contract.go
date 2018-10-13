package radius

type Radius interface {
	CreateCredentials(plan PricingPlan) (string, error)
}
