package admin

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/strusty/worldbit-wifi/database"
)

type store struct {
	db *gorm.DB
}

func NewAdminStore(db *gorm.DB) database.AdminStore {
	db.AutoMigrate(&database.Admin{})
	return store{
		db: db,
	}
}

func (store store) Create(admin *database.Admin) error {
	guid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	admin.ID = guid.String()
	return store.db.Create(admin).Error
}

func (store store) Update(id string, key string, value interface{}) error {
	return store.db.Model(&database.Admin{}).Where("id = ?", id).Update(key, value).Error
}

func (store store) ByLogin(login string) (*database.Admin, error) {
	admin := new(database.Admin)
	if err := store.db.Where("login = ?", login).First(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

func (store store) ByID(id string) (*database.Admin, error) {
	admin := new(database.Admin)
	if err := store.db.Where("id = ?", id).First(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}
