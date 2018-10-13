package pricing_plans

type PricingPlans interface {
	Create(plan *PricingPlan) error
	Update(plan *PricingPlan) error
	Delete(id string) error
	All() ([]PricingPlan, error)
	ByID(id string) (*PricingPlan, error)
}
