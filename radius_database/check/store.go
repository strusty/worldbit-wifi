package check

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database"
	"github.com/jinzhu/gorm"
)

type store struct {
	db *gorm.DB
}

func New(db *gorm.DB) radius_database.CheckStore {
	return store{
		db: db,
	}
}

func (store store) Create(check *radius_database.Check) error {
	return store.db.Create(check).Error
}
