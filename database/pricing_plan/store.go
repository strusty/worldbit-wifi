package pricing_plan

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"github.com/jinzhu/gorm"
)

type store struct {
	db *gorm.DB
}

func NewPricingPlanStore(db *gorm.DB) database.PricingPlanStore {
	db.AutoMigrate(&database.PricingPlan{})
	return store{
		db: db,
	}
}

func (store store) Create(plan *database.PricingPlan) error {
	return store.db.Create(plan).Error
}

func (store store) Update(plan *database.PricingPlan) error {
	return store.db.Save(plan).Error
}

func (store store) Delete(plan *database.PricingPlan) error {
	return store.db.Delete(plan).Error
}

func (store store) ByID(id string) (*database.PricingPlan, error) {
	entity := new(database.PricingPlan)

	if err := store.db.
		Where("id = ?", id).
		First(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (store store) All() ([]database.PricingPlan, error) {
	entity := new([]database.PricingPlan)

	if err := store.db.
		Find(entity).Error; err != nil {
		return nil, err
	}

	return *entity, nil
}
