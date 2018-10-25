package database

type AuthenticationsStore interface {
	Create(entity *Authentication) error
	ByPhoneNumber(phoneNumber string) (*Authentication, error)
	ByConfirmationCode(confirmationCode string) (*Authentication, error)
}

type PricingPlanStore interface {
	Create(plan *PricingPlan) error
	Update(plan *PricingPlan) error
	Delete(plan *PricingPlan) error
	All() ([]PricingPlan, error)
	ByID(id string) (*PricingPlan, error)
}

type AdminStore interface {
	Create(admin *Admin) error
	Update(id string, key string, value interface{}) error
	ByLogin(login string) (*Admin, error)
	ByID(id string) (*Admin, error)
}

type SalesStore interface {
	Create(sale *UsedSale) error
	ByPayPalSaleID(saleID string) (*UsedSale, error)
}
