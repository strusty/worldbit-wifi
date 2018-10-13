package pricing_plans

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"github.com/gofrs/uuid"
)

type service struct {
	store database.PricingPlanStore
}

func New(store database.PricingPlanStore) PricingPlans {
	return service{
		store: store,
	}
}

func (service service) Create(plan *PricingPlan) error {
	guid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	plan.ID = guid.String()
	return service.store.Create(&database.PricingPlan{
		ID:        plan.ID,
		AmountUSD: plan.AmountUSD,
		Duration:  plan.Duration,
		MaxUsers:  plan.MaxUsers,
		UpLimit:   plan.UpLimit,
		DownLimit: plan.DownLimit,
		PurgeDays: plan.PurgeDays,
	})
}

func (service service) Update(plan *PricingPlan) error {
	return service.store.Update(&database.PricingPlan{
		ID:        plan.ID,
		AmountUSD: plan.AmountUSD,
		Duration:  plan.Duration,
		MaxUsers:  plan.MaxUsers,
		UpLimit:   plan.UpLimit,
		DownLimit: plan.DownLimit,
		PurgeDays: plan.PurgeDays,
	})
}

func (service service) Delete(id string) error {
	return service.store.Delete(&database.PricingPlan{
		ID: id,
	})
}

func (service service) ByID(id string) (*PricingPlan, error) {
	pricingPlan, err := service.store.ByID(id)
	if err != nil {
		return nil, err
	}

	return &PricingPlan{
		ID:        pricingPlan.ID,
		AmountUSD: pricingPlan.AmountUSD,
		Duration:  pricingPlan.Duration,
		MaxUsers:  pricingPlan.MaxUsers,
		UpLimit:   pricingPlan.UpLimit,
		DownLimit: pricingPlan.DownLimit,
		PurgeDays: pricingPlan.PurgeDays,
	}, nil
}

func (service service) All() ([]PricingPlan, error) {
	pricingPlans, err := service.store.All()
	if err != nil {
		return nil, err
	}

	outputPricingPlans := make([]PricingPlan, len(pricingPlans))

	for i, pricingPlan := range pricingPlans {
		outputPricingPlans[i] = PricingPlan{
			ID:        pricingPlan.ID,
			AmountUSD: pricingPlan.AmountUSD,
			Duration:  pricingPlan.Duration,
			MaxUsers:  pricingPlan.MaxUsers,
			UpLimit:   pricingPlan.UpLimit,
			DownLimit: pricingPlan.DownLimit,
			PurgeDays: pricingPlan.PurgeDays,
		}
	}
	return outputPricingPlans, nil
}
