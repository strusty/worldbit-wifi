package check

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database"
	"github.com/jinzhu/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) radius_database.CheckStore {
	return store{
		db: db,
	}
}

func (store store) Create(check *radius_database.Check) error {
	return store.db.Create(check).Error
}

func (store store) SessionChecks() ([]radius_database.Check, error) {
	check := make([]radius_database.Check, 0)

	if err := store.db.
		Where("attribute = ?", "Max-Daily-Session").
		Find(&check).Error; err != nil {
		return nil, err
	}

	return check, nil
}

func (store store) DeleteChecksByUsername(username string) error {
	return store.db.Where("username = ?", username).Delete(&radius_database.Check{}).Error
}
